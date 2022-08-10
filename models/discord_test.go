package models

import (
	"path/filepath"
	"testing"
)

func TestUser(t *testing.T) {

	path := filepath.Join(t.TempDir(), "test.sqlite")
	ConnectDatabase(path)

	user := DiscordUser{
		ID:            "1234",
		Username:      "andrew",
		Discriminator: "6969",
		WordleGames:   nil,
		WordleStats: WordleStat{
			UserID:          "1234",
			AverageScore:    7,
			AverageFirstRow: 1,
			GamesPlayed:     69,
			PlayedToday:     false,
		},
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Errorf("Error creating DiscordUser with id %s", user.ID)
	}

	user = DiscordUser{}
	if err := DB.Where("id = ?", "1234").First(&user).Error; err != nil {
		t.Errorf("Error reading DiscordUser with id %s", user.ID)
	}
}
