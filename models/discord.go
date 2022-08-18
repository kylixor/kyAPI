package models

import (
	"time"

	"gorm.io/gorm"
)

type DiscordUser struct {
	ID                 string `json:"id" binding:"required" gorm:"primaryKey"` // Discord User ID
	Username           string `json:"username"`
	Discriminator      string `json:"discriminator"`
	GetWordleReminders bool
	WordleGames        []WordleGame `gorm:"foreignKey:UserID;references:ID"`
	WordleStats        WordleStat   `gorm:"foreignKey:UserID;references:ID"`
}

func GetOrCreateDiscordUser(id string) (*DiscordUser, error) {
	var user DiscordUser
	err := DB.Take(&user, id).Error
	if err == nil {
		return &user, nil
	} else if err == gorm.ErrRecordNotFound {
		user := DiscordUser{
			ID: id,
		}
		// ENHANCEMENT: Lookup user in Discord?
		if err := DB.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, err
	}
}

func (u *DiscordUser) CalculateWordleStats() (err error) {
	var ws WordleStat
	if err := DB.Take(&ws, u.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ws.UserID = u.ID
			if err := DB.Create(&ws).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Gather and Sum all Wordle Scores
	var wordles []WordleGame
	if err := DB.Find(&wordles, &WordleGame{UserID: u.ID}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ws.AverageFirstRow = 0
			ws.AverageScore = 0
			ws.GamesPlayed = 0
			ws.PlayedToday = false
			DB.Save(&ws)
			return nil
		} else {
			return err
		}
	}

	var sum float32
	today_wordle := uint16(time.Since(WORDLE_DAY_0).Hours() / 24)
	for _, wordle := range wordles {
		sum += float32(wordle.Score)
		if wordle.Day == today_wordle {
			ws.PlayedToday = true
		}
	}
	ws.AverageFirstRow = 0
	ws.AverageScore = sum / float32(len(wordles))
	ws.GamesPlayed = uint16(len(wordles))
	DB.Save(&ws)
	return nil
}
