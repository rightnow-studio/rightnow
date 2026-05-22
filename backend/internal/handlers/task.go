package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cike/internal/models"
)

type CreateTaskReq struct {
	TextKey         string `json:"text_key" binding:"required"`
	DurationSeconds int    `json:"duration_seconds" binding:"required,min=10,max=300"`
	Category        string `json:"category" binding:"required"`
}

func ListTasks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		rows, err := db.Query(`
			SELECT t.id, t.text_key, t.duration_seconds, t.category, t.is_preset, t.created_by, t.is_active, t.created_at
			FROM tasks t
			LEFT JOIN user_task_pool utp ON t.id = utp.task_id AND utp.user_id = ?
			WHERE t.is_preset = TRUE OR utp.user_id = ?
			ORDER BY t.created_at DESC
		`, userID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tasks []models.Task
		for rows.Next() {
			var t models.Task
			if err := rows.Scan(&t.ID, &t.TextKey, &t.DurationSeconds, &t.Category, &t.IsPreset, &t.CreatedBy, &t.IsActive, &t.CreatedAt); err != nil {
				continue
			}
			tasks = append(tasks, t)
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

func CreateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var req CreateTaskReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := uuid.NewString()
		now := time.Now()
		_, err := db.Exec(
			"INSERT INTO tasks (id, text_key, duration_seconds, category, is_preset, created_by, is_active, created_at) VALUES (?, ?, ?, ?, FALSE, ?, TRUE, ?)",
			id, req.TextKey, req.DurationSeconds, req.Category, userID, now,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = db.Exec("INSERT INTO user_task_pool (user_id, task_id, added_at) VALUES (?, ?, ?)", userID, id, now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func UpdateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		taskID := c.Param("id")
		var req struct {
			IsActive *bool `json:"is_active"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.IsActive != nil {
			_, err := db.Exec("UPDATE tasks SET is_active = ? WHERE id = ? AND created_by = ?", *req.IsActive, taskID, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func DeleteTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		taskID := c.Param("id")
		_, err := db.Exec("DELETE FROM tasks WHERE id = ? AND created_by = ? AND is_preset = FALSE", taskID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, _ = db.Exec("DELETE FROM user_task_pool WHERE task_id = ?", taskID)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
