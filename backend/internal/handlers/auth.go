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
