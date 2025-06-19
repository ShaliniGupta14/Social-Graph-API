package main

import (
	"social_graph_api/db"
	"social_graph_api/router"
)

func main() {
	db.ConnectDB()
	r := router.SetupRouter()
	r.Run(":8080")
}
