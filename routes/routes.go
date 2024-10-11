package routes

import (
    "github.com/gin-gonic/gin"
    "cassandra_go_recommendation_api/controllers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.POST("/items", controllers.CreateItem)
    r.POST("/interactions", controllers.LogInteraction)
    r.GET("/recommendations/:user_id", controllers.GetRecommendations)

    return r
}
