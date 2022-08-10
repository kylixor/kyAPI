package models

import "time"

type WordleGame struct {
	MessageID string `gorm:"primaryKey"`
	UserID    string
	ChannelID string
	Day       uint16 // Wordle day number
	Score     uint8  // Score out of 6
	Row1      string
	Row1Score int8
	Row2      string
	Row2Score int8
	Row3      string
	Row3Score int8
	Row4      string
	Row4Score int8
	Row5      string
	Row5Score int8
	Row6      string
	Row6Score int8
}

type WordleStat struct {
	UserID          string `gorm:"primaryKey"`
	AverageScore    float32
	AverageFirstRow float32
	GamesPlayed     uint16
	PlayedToday     bool
}

var WORDLE_DAY_0 = time.Date(2021, time.June, 19, 0, 0, 0, 0, time.Now().Location())
