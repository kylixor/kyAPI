package controllers

import (
	"kyapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindDiscordUser(c *gin.Context) {
	var user models.DiscordUser
	if err := models.DB.Where("id = ?", c.Param("id")).First(user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})
		return
	}
	c.JSON(http.StatusAccepted, &user)
}

func FindDiscordUsers(c *gin.Context) {
	var users []models.DiscordUser
	models.DB.Find(&users)
}

func CreateDiscordUser(c *gin.Context) {
	var user models.DiscordUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.DiscordUser{
		ID:                 user.ID,
		Username:           "",
		Discriminator:      "",
		WordleGames:        nil,
		GetWordleReminders: false,
	}

	if err := models.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := newUser.CalculateWordleStats(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func UpdateDiscordUser(c *gin.Context) {
	var user models.DiscordUser

	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	var updateUser models.DiscordUser

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Model(&user).Updates(models.DiscordUser{Username: updateUser.Username, Discriminator: updateUser.Discriminator}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
