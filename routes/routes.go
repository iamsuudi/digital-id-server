package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/iamsuudi/digital-id-server/controllers"
)

func RegisterRoutes(r *gin.Engine) {
    userGroup := r.Group("/users")
    {
        userGroup.GET("/", controllers.GetAllUsers)
        userGroup.GET("/:id", controllers.GetUserByID)
        userGroup.POST("/", controllers.CreateUser)
    }
}
