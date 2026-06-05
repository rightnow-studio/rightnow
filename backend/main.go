package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"cike/internal/config"
	"cike/internal/db"
	"cike/internal/handlers"
	"cike/internal/middleware"
)

func main() {
	cfg := config.Load()

	database, err := db.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	if err := db.SeedAdmin(database, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.POST("/auth/register", handlers.Register(database))
		api.POST("/auth/login", handlers.Login(database))
		api.POST("/auth/refresh", handlers.RefreshToken(database))
		api.POST("/auth/guest", handlers.CreateGuest(database))

		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			authorized.GET("/me", handlers.GetMe(database))
			authorized.POST("/auth/bind", handlers.BindEmail(database))
			authorized.GET("/dispatch", handlers.DispatchTask(database))
			authorized.GET("/tasks", handlers.ListTasks(database))
			authorized.POST("/tasks", handlers.CreateTask(database))
			authorized.PATCH("/tasks/:id", handlers.UpdateTask(database))
			authorized.DELETE("/tasks/:id", handlers.DeleteTask(database))
			authorized.GET("/history", handlers.GetHistory(database))
			authorized.GET("/history/stats", handlers.GetStats(database))
			authorized.POST("/history", handlers.RecordHistory(database))
			authorized.GET("/settings", handlers.GetSettings(database))
			authorized.PUT("/settings", handlers.UpdateSettings(database))
		}
	}

	admin := r.Group("/admin")
	{
		admin.POST("/auth/login", handlers.AdminLogin(database, cfg.AdminSecret))

		adminAuth := admin.Group("/")
		adminAuth.Use(middleware.AdminJWTAuth(cfg.AdminSecret))
		{
			adminAuth.GET("/tasks", handlers.AdminListTasks(database))
			adminAuth.POST("/tasks", handlers.AdminCreateTask(database))
			adminAuth.PUT("/tasks/:id", handlers.AdminUpdateTask(database))
			adminAuth.DELETE("/tasks/:id", handlers.AdminDeleteTask(database))
			adminAuth.GET("/categories", handlers.AdminListCategories(database))
			adminAuth.POST("/categories", handlers.AdminCreateCategory(database))
			adminAuth.DELETE("/categories/:id", handlers.AdminDeleteCategory(database))
		}
	}

	log.Printf("server running on %s", cfg.ServerAddr)
	if err := r.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
