package handlers

import (
	"net/http"

	"social_graph_api/db"
	"social_graph_api/models"

	"github.com/gin-gonic/gin"
)

type ConnectRequest struct {
	UserID   uint `json:"user_id"`
	TargetID uint `json:"target_id"`
}

func ConnectUsers(c *gin.Context) {
	var req ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	var target models.User

	if err := db.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.DB.First(&target, req.TargetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	if err := db.DB.Model(&user).Association("Connections").Append(&target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect users"})
		return
	}

	if err := db.DB.Model(&target).Association("Connections").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect users back"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Users connected successfully!"})
}

func GetConnections(c *gin.Context) {
	var user models.User
	if err := db.DB.Preload("Connections").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user.Connections)
}

func GetRecommendations(c *gin.Context) {
	var user models.User
	if err := db.DB.Preload("Connections").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 1st-degree connections map
	firstDegree := make(map[uint]bool)
	for _, conn := range user.Connections {
		firstDegree[conn.ID] = true
	}

	recommendations := []models.User{}
	for _, conn := range user.Connections {
		var secondDegreeUser models.User
		if err := db.DB.Preload("Connections").First(&secondDegreeUser, conn.ID).Error; err != nil {
			continue
		}
		for _, friendOfFriend := range secondDegreeUser.Connections {
			// Skip if self or already a connection
			if friendOfFriend.ID == user.ID || firstDegree[friendOfFriend.ID] {
				continue
			}
			// Avoid duplicates
			exists := false
			for _, rec := range recommendations {
				if rec.ID == friendOfFriend.ID {
					exists = true
					break
				}
			}
			if !exists {
				recommendations = append(recommendations, *friendOfFriend)
			}
		}
	}

	c.JSON(http.StatusOK, recommendations)
}
