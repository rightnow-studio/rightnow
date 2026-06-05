package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"cike/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin POST /admin/auth/login
func AdminLogin(db *sql.DB, adminSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
			return
		}

		var id, hash string
		err := db.QueryRow(
			"SELECT id, password_hash FROM admin_users WHERE username = ?",
			req.Username,
		).Scan(&id, &hash)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := middleware.GenerateAdminToken(id, adminSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":    token,
			"username": req.Username,
		})
	}
}

// AdminListTasks GET /admin/tasks
func AdminListTasks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(
			"SELECT id, text_key, duration_seconds, category, is_preset, is_active, created_at FROM tasks WHERE is_preset = TRUE ORDER BY id ASC",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		defer rows.Close()

		type TaskItem struct {
			ID              string    `json:"id"`
			TextKey         string    `json:"text_key"`
			DurationSeconds int       `json:"duration_seconds"`
			Category        string    `json:"category"`
			IsPreset        bool      `json:"is_preset"`
			IsActive        bool      `json:"is_active"`
			CreatedAt       time.Time `json:"created_at"`
		}

		tasks := []TaskItem{}
		for rows.Next() {
			var t TaskItem
			if err := rows.Scan(&t.ID, &t.TextKey, &t.DurationSeconds, &t.Category, &t.IsPreset, &t.IsActive, &t.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
				return
			}
			tasks = append(tasks, t)
		}
		c.JSON(http.StatusOK, tasks)
	}
}

// AdminCreateTask POST /admin/tasks
func AdminCreateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TextKey         string `json:"text_key" binding:"required"`
			DurationSeconds int    `json:"duration_seconds" binding:"required,min=10,max=300"`
			Category        string `json:"category" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证分类存在
		var catCount int
		db.QueryRow("SELECT COUNT(*) FROM task_categories WHERE name = ? AND is_active = TRUE", req.Category).Scan(&catCount)
		if catCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
			return
		}

		id := uuid.New().String()
		_, err := db.Exec(
			"INSERT INTO tasks (id, text_key, duration_seconds, category, is_preset, is_active) VALUES (?, ?, ?, ?, TRUE, TRUE)",
			id, req.TextKey, req.DurationSeconds, req.Category,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

// AdminUpdateTask PUT /admin/tasks/:id
func AdminUpdateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")
		var req struct {
			TextKey         *string `json:"text_key"`
			DurationSeconds *int    `json:"duration_seconds"`
			Category        *string `json:"category"`
			IsActive        *bool   `json:"is_active"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 确认是预设任务
		var isPreset bool
		err := db.QueryRow("SELECT is_preset FROM tasks WHERE id = ?", taskID).Scan(&isPreset)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		if !isPreset {
			c.JSON(http.StatusForbidden, gin.H{"error": "only preset tasks can be edited via admin"})
			return
		}

		if req.TextKey != nil {
			db.Exec("UPDATE tasks SET text_key = ? WHERE id = ?", *req.TextKey, taskID)
		}
		if req.DurationSeconds != nil {
			if *req.DurationSeconds < 10 || *req.DurationSeconds > 300 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "duration must be 10-300"})
				return
			}
			db.Exec("UPDATE tasks SET duration_seconds = ? WHERE id = ?", *req.DurationSeconds, taskID)
		}
		if req.Category != nil {
			var catCount int
			db.QueryRow("SELECT COUNT(*) FROM task_categories WHERE name = ? AND is_active = TRUE", *req.Category).Scan(&catCount)
			if catCount == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
				return
			}
			db.Exec("UPDATE tasks SET category = ? WHERE id = ?", *req.Category, taskID)
		}
		if req.IsActive != nil {
			db.Exec("UPDATE tasks SET is_active = ? WHERE id = ?", *req.IsActive, taskID)
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

// AdminDeleteTask DELETE /admin/tasks/:id
func AdminDeleteTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")

		var isPreset bool
		err := db.QueryRow("SELECT is_preset FROM tasks WHERE id = ?", taskID).Scan(&isPreset)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		if !isPreset {
			c.JSON(http.StatusForbidden, gin.H{"error": "only preset tasks can be deleted via admin"})
			return
		}

		db.Exec("DELETE FROM user_task_pool WHERE task_id = ?", taskID)
		db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

// AdminListCategories GET /admin/categories
func AdminListCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, label_key, is_active FROM task_categories ORDER BY name ASC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		defer rows.Close()

		type CatItem struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			LabelKey string `json:"label_key"`
			IsActive bool   `json:"is_active"`
		}

		cats := []CatItem{}
		for rows.Next() {
			var cat CatItem
			if err := rows.Scan(&cat.ID, &cat.Name, &cat.LabelKey, &cat.IsActive); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
				return
			}
			cats = append(cats, cat)
		}
		c.JSON(http.StatusOK, cats)
	}
}

// AdminCreateCategory POST /admin/categories
func AdminCreateCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			LabelKey string `json:"label_key" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := "cat_" + uuid.New().String()[:8]
		_, err := db.Exec(
			"INSERT INTO task_categories (id, name, label_key) VALUES (?, ?, ?)",
			id, req.Name, req.LabelKey,
		)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

// AdminDeleteCategory DELETE /admin/categories/:id
func AdminDeleteCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catID := c.Param("id")

		// 检查是否有任务引用该分类
		var catName string
		err := db.QueryRow("SELECT name FROM task_categories WHERE id = ?", catID).Scan(&catName)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}

		var taskCount int
		db.QueryRow("SELECT COUNT(*) FROM tasks WHERE category = ?", catName).Scan(&taskCount)
		if taskCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "category has tasks, cannot delete", "task_count": taskCount})
			return
		}

		db.Exec("DELETE FROM task_categories WHERE id = ?", catID)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
