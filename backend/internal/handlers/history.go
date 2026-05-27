package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecordHistoryReq struct {
	TaskID  string `json:"task_id" binding:"required"`
	Outcome string `json:"outcome" binding:"required,oneof=completed skipped"`
}

type HistoryItem struct {
	ID        string         `json:"id"`
	TaskID    string         `json:"task_id"`
	TextKey   string         `json:"text_key"`
	Category  string         `json:"category"`
	Duration  int            `json:"duration_seconds"`
	Outcome   string         `json:"outcome"`
	CreatedAt time.Time      `json:"created_at"`
}

func RecordHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var req RecordHistoryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := uuid.NewString()
		now := time.Now()
		_, err := db.Exec(
			"INSERT INTO task_history (id, user_id, task_id, started_at, ended_at, outcome, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id, userID, req.TaskID, now, now, req.Outcome, now,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func GetHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		rows, err := db.Query(
			`SELECT th.id, th.task_id, t.text_key, t.category, t.duration_seconds, th.outcome, th.created_at
			 FROM task_history th
			 JOIN tasks t ON th.task_id = t.id
			 WHERE th.user_id = ?
			 ORDER BY th.created_at DESC`,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var items []HistoryItem
		for rows.Next() {
			var item HistoryItem
			if err := rows.Scan(&item.ID, &item.TaskID, &item.TextKey, &item.Category, &item.Duration, &item.Outcome, &item.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			items = append(items, item)
		}
		c.JSON(http.StatusOK, items)
	}
}

func GetStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var totalCompleted, totalSkipped int
		err := db.QueryRow("SELECT COUNT(*) FROM task_history WHERE user_id = ? AND outcome = 'completed'", userID).Scan(&totalCompleted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM task_history WHERE user_id = ? AND outcome = 'skipped'", userID).Scan(&totalSkipped)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"total_completed": totalCompleted,
			"total_skipped":   totalSkipped,
		})
	}
}
