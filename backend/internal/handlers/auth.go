package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"cike/internal/config"
	"cike/internal/middleware"
	"cike/internal/models"
)

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		var req RegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		id := uuid.NewString()
		_, err = db.Exec("INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)", id, req.Email, string(hash))
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		access, refresh, err := middleware.GenerateTokens(id, cfg.JWTSecret, cfg.RefreshSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
			return
		}
		c.JSON(http.StatusCreated, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		var req LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		var hash string
		err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", req.Email).Scan(&user.ID, &user.Email, &hash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		access, refresh, err := middleware.GenerateTokens(user.ID, cfg.JWTSecret, cfg.RefreshSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
			return
		}
		c.JSON(http.StatusOK, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
	}
}

func GetMe(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var user models.User
		err := db.QueryRow("SELECT id, email, is_temp, created_at FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.IsTemp, &user.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func RefreshToken(db *sql.DB) gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims, err := middleware.ParseRefreshToken(req.RefreshToken, cfg.RefreshSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		access, refresh, err := middleware.GenerateTokens(userID, cfg.JWTSecret, cfg.RefreshSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
			return
		}
		c.JSON(http.StatusOK, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
	}
}

func CreateGuest(db *sql.DB) gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		id := uuid.NewString()
		_, err := db.Exec("INSERT INTO users (id, email, password_hash, is_temp) VALUES (?, '', '', TRUE)", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest"})
			return
		}
		access, refresh, err := middleware.GenerateTokens(id, cfg.JWTSecret, cfg.RefreshSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
			return
		}
		c.JSON(http.StatusCreated, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
	}
}

func BindEmail(db *sql.DB) gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		guestID := c.GetString("user_id")

		var isTemp bool
		if err := db.QueryRow("SELECT is_temp FROM users WHERE id = ?", guestID).Scan(&isTemp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !isTemp {
			c.JSON(http.StatusForbidden, gin.H{"error": "only guest users can bind"})
			return
		}

		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingID, existingHash string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", req.Email).Scan(&existingID, &existingHash)

		if err == sql.ErrNoRows {
			// 全新邮箱：升级临时用户
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}
			if _, err := db.Exec("UPDATE users SET email = ?, password_hash = ?, is_temp = FALSE WHERE id = ?", req.Email, string(hash), guestID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			access, refresh, _ := middleware.GenerateTokens(guestID, cfg.JWTSecret, cfg.RefreshSecret)
			c.JSON(http.StatusOK, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 邮箱已存在：验证密码后合并数据
		if err := bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if _, err := tx.Exec("UPDATE task_history SET user_id = ? WHERE user_id = ?", existingID, guestID); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if _, err := tx.Exec("DELETE FROM users WHERE id = ?", guestID); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		access, refresh, _ := middleware.GenerateTokens(existingID, cfg.JWTSecret, cfg.RefreshSecret)
		c.JSON(http.StatusOK, TokenResp{AccessToken: access, RefreshToken: refresh, ExpiresIn: 900})
	}
}
