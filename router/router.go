package router

import (
	"social_graph_api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetAllUsers)

	r.POST("/connect", handlers.ConnectUsers)
	r.GET("/connections/:id", handlers.GetConnections)
	r.GET("/recommendations/:id", handlers.GetRecommendations)

	return r
}
