package utilities

import (
	"net/http"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/gin-gonic/gin"
)

func GetLoggedInUser(c *gin.Context) *models.User {
	val, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil
	}
	user := val.(*models.User)
	return user
}
