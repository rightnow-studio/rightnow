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
