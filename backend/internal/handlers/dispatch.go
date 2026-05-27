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

		// 查询用户个人任务（非预设）
		customQuery := `
			SELECT t.id, t.text_key, t.duration_seconds, t.category
			FROM tasks t
			JOIN user_task_pool utp ON t.id = utp.task_id AND utp.user_id = ?
			WHERE t.is_preset = FALSE AND t.is_active = TRUE
		`
		customArgs := []interface{}{userID}
		for _, rid := range recentIDs {
			customQuery += " AND t.id != ?"
			customArgs = append(customArgs, rid)
		}

		// 查询系统预设任务
		presetQuery := `
			SELECT t.id, t.text_key, t.duration_seconds, t.category
			FROM tasks t
			WHERE t.is_preset = TRUE AND t.is_active = TRUE
		`
		presetArgs := []interface{}{}
		for _, rid := range recentIDs {
			presetQuery += " AND t.id != ?"
			presetArgs = append(presetArgs, rid)
		}

		var customCandidates []DispatchResp
		customRows, err := db.Query(customQuery, customArgs...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for customRows.Next() {
			var d DispatchResp
			if err := customRows.Scan(&d.ID, &d.TextKey, &d.DurationSeconds, &d.Category); err == nil {
				customCandidates = append(customCandidates, d)
			}
		}
		customRows.Close()

		var presetCandidates []DispatchResp
		presetRows, err := db.Query(presetQuery, presetArgs...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for presetRows.Next() {
			var d DispatchResp
			if err := presetRows.Scan(&d.ID, &d.TextKey, &d.DurationSeconds, &d.Category); err == nil {
				presetCandidates = append(presetCandidates, d)
			}
		}
		presetRows.Close()

		// 读取用户自定义概率，默认 60%
		var customRatio int
		err = db.QueryRow("SELECT custom_task_ratio FROM user_settings WHERE user_id = ?", userID).Scan(&customRatio)
		if err == sql.ErrNoRows {
			customRatio = 60
		}
		rand.Seed(time.Now().UnixNano())
		useCustom := rand.Float64() < float64(customRatio)/100.0

		var candidates []DispatchResp
		if useCustom && len(customCandidates) > 0 {
			candidates = customCandidates
		} else if len(presetCandidates) > 0 {
			candidates = presetCandidates
		} else if len(customCandidates) > 0 {
			candidates = customCandidates
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

		selected := candidates[rand.Intn(len(candidates))]
		c.JSON(http.StatusOK, selected)
	}
}
