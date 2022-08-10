package controllers

import (
	"kyapi/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewWordleGame struct {
	MessageID string `json:"message_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
	Day       uint16 `json:"day" binding:"required"`
	Score     uint8  `json:"score" binding:"required"`
	Row1      string `json:"row1" binding:"required"`
	Row2      string `json:"row2"`
	Row3      string `json:"row3"`
	Row4      string `json:"row4"`
	Row5      string `json:"row5"`
	Row6      string `json:"row6"`
}

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
	var wordle NewWordleGame
	var newWordle models.WordleGame

	if err := c.ShouldBindJSON(&wordle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Take(&newWordle, wordle.MessageID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "wordle game already exits"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newWordle = models.WordleGame{
		MessageID: wordle.MessageID,
		UserID:    wordle.UserID,
		ChannelID: wordle.ChannelID,
		Day:       wordle.Day,
		Score:     wordle.Score,
		Row1:      wordle.Row1,
		Row1Score: rowScoreFromRow(wordle.Row1),
		Row2:      wordle.Row2,
		Row2Score: rowScoreFromRow(wordle.Row2),
		Row3:      wordle.Row3,
		Row3Score: rowScoreFromRow(wordle.Row3),
		Row4:      wordle.Row4,
		Row4Score: rowScoreFromRow(wordle.Row4),
		Row5:      wordle.Row5,
		Row5Score: rowScoreFromRow(wordle.Row5),
		Row6:      wordle.Row6,
		Row6Score: rowScoreFromRow(wordle.Row6),
	}

	if err := models.DB.Create(&newWordle).Error; err != nil {
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
	user.WordleGames = append(user.WordleGames, newWordle)
	models.DB.Save(&user)

	c.JSON(http.StatusAccepted, newWordle)
}
