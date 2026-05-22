package handlers

import (
	"database/sql"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DispatchResp struct {
	ID              string `json:"id"`
	TextKey         string `json:"text_key"`
	DurationSeconds int    `json:"duration_seconds"`
	Category        string `json:"category"`
}

func DispatchTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		var recentIDs []string
		rows, err := db.Query(
			"SELECT task_id FROM task_history WHERE user_id = ? ORDER BY started_at DESC LIMIT 5",
			userID,
		)
		if err == nil {
			for rows.Next() {
				var tid string
				if err := rows.Scan(&tid); err == nil {
					recentIDs = append(recentIDs, tid)
				}
			}
			rows.Close()
		}

		query := `
			SELECT t.id, t.text_key, t.duration_seconds, t.category
			FROM tasks t
			LEFT JOIN user_task_pool utp ON t.id = utp.task_id AND utp.user_id = ?
			WHERE (t.is_preset = TRUE OR utp.user_id = ?)
			  AND t.is_active = TRUE
		`
		args := []interface{}{userID, userID}
		for _, rid := range recentIDs {
			query += " AND t.id != ?"
			args = append(args, rid)
		}

		rows, err = db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var candidates []DispatchResp
		for rows.Next() {
			var d DispatchResp
			if err := rows.Scan(&d.ID, &d.TextKey, &d.DurationSeconds, &d.Category); err == nil {
				candidates = append(candidates, d)
			}
		}

		if len(candidates) == 0 {
			c.JSON(http.StatusOK, DispatchResp{
				ID:              "fallback",
				TextKey:         "task.write_start_on_paper",
				DurationSeconds: 30,
				Category:        "admin",
			})
			return
		}

		rand.Seed(time.Now().UnixNano())
		selected := candidates[rand.Intn(len(candidates))]
		c.JSON(http.StatusOK, selected)
	}
}
