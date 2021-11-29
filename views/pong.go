package views

import (
	"net/http"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	"github.com/gin-gonic/gin"
)

// Pong test server is operating
func Pong(c *gin.Context) {
	if models.Ping() != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
