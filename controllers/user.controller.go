package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iamsuudi/digital-id-server/models"
	"github.com/iamsuudi/digital-id-server/services"
)

func GetAllUsers(c *gin.Context) {
    users, err := services.GetUsers()
    fmt.Println(err)
    c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    user, err := services.GetUserByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    services.CreateUser(user)
    c.JSON(http.StatusCreated, user)
}
