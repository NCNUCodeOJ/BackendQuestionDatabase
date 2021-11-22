package router

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/views"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
		if err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "系統錯誤",
				"error":   err.Error(),
			})
		} else {
			c.Set("userID", uint(id))
			c.Next()
		}
	}
}

// SetupRouter index
func SetupRouter() *gin.Engine {
	if gin.Mode() == "test" {
		err := godotenv.Load(".env.test")
		if err != nil {
			log.Println("Error loading .env file")
		}
	} else if gin.Mode() == "debug" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "NCNUOJ",
		SigningAlgorithm: "HS512",
		Key:              []byte(os.Getenv("SECRET_KEY")),
		MaxRefresh:       time.Hour,
		TimeFunc:         time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	baseURL := "api/v1"
	r := gin.Default()
	r.GET("/ping", views.Pong)

	problem := r.Group(baseURL + "/problem")
	problem.Use(authMiddleware.MiddlewareFunc())
	problem.Use(getUserID())
	{
		problem.POST("", views.CreateProblem) // 創建題目
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}
