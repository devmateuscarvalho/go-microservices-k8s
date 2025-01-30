package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key"))

type Task struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserID    uint   `json:"user_id"`
}

func main() {
	// Configuração do banco de dados
	dsn := getEnvOrDefault("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=microservices port=5432 sslmode=disable")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	db.AutoMigrate(&Task{})

	// Configuração do Gin
	r := gin.Default()

	// Middleware de autenticação
	authorized := r.Group("/")
	authorized.Use(authMiddleware())
	{
		authorized.GET("/tasks", getTasks)
		authorized.POST("/tasks", createTask)
		authorized.PUT("/tasks/:id", updateTask)
		authorized.DELETE("/tasks/:id", deleteTask)
	}

	// Rota de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Iniciar servidor
	port := getEnvOrDefault("PORT", "8081")
	r.Run(":" + port)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("user_id", uint(claims["sub"].(float64)))
		c.Next()
	}
}

func getTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var tasks []Task
	db.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(200, tasks)
}

func createTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID
	db.Create(&task)
	c.JSON(201, task)
}

func updateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("id")

	var task Task
	if err := db.First(&task, taskID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	if task.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	var updateData Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.Title = updateData.Title
	task.Completed = updateData.Completed
	db.Save(&task)

	c.JSON(200, task)
}

func deleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("id")

	var task Task
	if err := db.First(&task, taskID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	if task.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	db.Delete(&task)
	c.Status(http.StatusNoContent)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
