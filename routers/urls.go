package router

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/views"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
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

	privateURL := "api/private/v1"
	r := gin.Default()

	// CORS
	if os.Getenv("FrontendURL") != "" {
		origins := strings.Split(os.Getenv("FrontendURL"), ",")
		log.Println("CORS:", origins)
		r.Use(cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
			AllowHeaders:     []string{"Origin, Authorization, Content-Type, Accept"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	r.GET("/ping", views.Pong)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	problem := r.Group(privateURL + "/problem")
	problem.Use(authMiddleware.MiddlewareFunc())
	problem.Use(getUserID())
	{
		// problem.GET("/tag/:tagName", views.GetProblemsByTag) // 查詢 該 tag 所有 problems
		problem.POST("", views.CreateProblem)                      // 創建題目
		problem.GET("/:id", views.GetProblemByID)                  // 取得題目
		problem.PATCH("/:id", views.EditProblem)                   // 編輯題目
		problem.POST("/:id/testcase", views.UploadProblemTestCase) // 上傳題目測試 test case

	}
	privateProblem := r.Group(privateURL + "/problem")
	privateProblem.Use(authMiddleware.MiddlewareFunc())
	privateProblem.Use(getUserID())
	{
		privateProblem.POST("/:id/submission", views.CreateSubmission) // 上傳 submission
	}
	submission := r.Group(privateURL + "/submission")
	{
		submission.PATCH("/:id/judge", views.UpdateSubmissionJudgeResult) // 更新 submission judge result
		submission.PATCH("/:id/style", views.UpdateSubmissionStyleResult) // 更新 submission style result
		submission.GET("/:id", views.GetSubmissionByID)                   // 取得 submission
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}
