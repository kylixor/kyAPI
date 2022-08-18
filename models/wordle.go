package models

import "time"

type WordleGame struct {
	MessageID string `json:"message_id" binding:"required" gorm:"primaryKey"`
	UserID    string `json:"user_id" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
	Day       uint16 `json:"day" binding:"required"`
	Score     uint8  `json:"score" binding:"required"`
	Row1      string `json:"row1" binding:"required"`
	Row1Score int8
	Row2      string `json:"row2"`
	Row2Score int8
	Row3      string `json:"row3"`
	Row3Score int8
	Row4      string `json:"row4"`
	Row4Score int8
	Row5      string `json:"row5"`
	Row5Score int8
	Row6      string `json:"row6"`
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
