package controllers

import (
	"time"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/btdjangbah001/chat-app/utilities"
	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	var groupInput models.CreateGroup

	if err := c.ShouldBindJSON(&groupInput); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	group := models.Group{
		Name:      groupInput.Name,
		OwnerID:   groupInput.OwnerID,
		CreatedAt: time.Now().UTC(),
	}

	if err := group.CreateGroup(); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := models.AddParticipantsWithUsernamesToGroup(group.ID, groupInput.Usernames); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to add participants to group"})
		return
	}

	c.JSON(200, gin.H{"group": group})
}

func GetGroupsForUser(c *gin.Context) {
	user := utilities.GetLoggedInUser(c)

	groups, err := models.GetGroupsForUser(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"groups": groups})
}
