package services

import (
	"errors"

	"github.com/iamsuudi/digital-id-server/config"
	"github.com/iamsuudi/digital-id-server/models"
)

func GetUsers() ([]models.User, error) {
    var users []models.User
    result := config.DB.Find(&users)
    return users, result.Error
}

func GetUserByID(id string) (models.User, error) {
    var user models.User
    result := config.DB.First(&user, id)
    if result.Error != nil {
        return user, errors.New("user not found")
    }
    return user, nil
}

func CreateUser(user models.User) error {
    return config.DB.Create(&user).Error
}
