package routes

import (
    "github.com/gin-gonic/gin"
    "cassandra_go_recommendation_api/controllers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.POST("/interactions", controllers.LogInteraction)
    r.GET("/recommendations/:user_id", controllers.GetRecommendations)
    r.GET("/popular-items", controllers.GetPopularItems)

    return r
}
