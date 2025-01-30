package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var db *gorm.DB
var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key"))

func main() {
	// Configuração do banco de dados
	dsn := getEnvOrDefault("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=microservices port=5432 sslmode=disable")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	db.AutoMigrate(&User{})

	// Configuração do Gin
	r := gin.Default()

	// Rotas
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Iniciar servidor
	port := getEnvOrDefault("PORT", "8080")
	r.Run(":" + port)
}

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Criar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
