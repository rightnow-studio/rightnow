package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateSettingsReq struct {
	CustomTaskRatio int `json:"custom_task_ratio" binding:"min=0,max=100"`
}

func GetSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var ratio int
		err := db.QueryRow("SELECT custom_task_ratio FROM user_settings WHERE user_id = ?", userID).Scan(&ratio)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"custom_task_ratio": 60})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"custom_task_ratio": ratio})
	}
}

func UpdateSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var req UpdateSettingsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := db.Exec(
			"INSERT INTO user_settings (user_id, custom_task_ratio, updated_at) VALUES (?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET custom_task_ratio = excluded.custom_task_ratio, updated_at = excluded.updated_at",
			userID, req.CustomTaskRatio, time.Now(),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"custom_task_ratio": req.CustomTaskRatio})
	}
}
