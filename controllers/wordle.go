package controllers

import (
	"kyapi/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func rowScoreFromRow(row string) (score int8) {
	if row == "" {
		return -1
	}
	yellows := int8(strings.Count(row, "🟨"))
	greens := int8(strings.Count(row, "🟩"))
	score = greens*2 + yellows*1
	return score
}

func CreateWordleGame(c *gin.Context) {
	var wordle models.WordleGame

	if err := c.ShouldBindJSON(&wordle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Take(&wordle, wordle.MessageID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "wordle game already exits"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	wordle.Row1Score = rowScoreFromRow(wordle.Row1)
	wordle.Row2Score = rowScoreFromRow(wordle.Row2)
	wordle.Row3Score = rowScoreFromRow(wordle.Row3)
	wordle.Row4Score = rowScoreFromRow(wordle.Row4)
	wordle.Row5Score = rowScoreFromRow(wordle.Row5)
	wordle.Row6Score = rowScoreFromRow(wordle.Row6)

	if err := models.DB.Create(&wordle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetOrCreateDiscordUser(wordle.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := user.CalculateWordleStats(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.WordleGames = append(user.WordleGames, wordle)
	models.DB.Save(&user)

	c.JSON(http.StatusAccepted, wordle)
}
