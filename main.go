package main

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var WORDLE_DAY_0 = time.Date(2021, time.June, 19, 0, 0, 0, 0, time.Now().Location())

type BarData struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	Data            []int  `json:"data"`
}

type GraphData struct {
	Labels []string  `json:"labels"`
	Bars   []BarData `json:"datasets"`
}

type WordleData struct {
	ID             string
	MessageID      string
	UserID         string
	ChannelID      string
	Day            int
	Score          int
	BlankCount     string
	YellowCount    string
	GreenCount     string
	FirstWordScore string
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func fakeData(c *gin.Context) {

	days_since_wordle_began := uint16(time.Since(WORDLE_DAY_0).Hours() / 24)

	calendar := make([]string, days_since_wordle_began)
	for day_index := range calendar {
		calendar[day_index] = WORDLE_DAY_0.AddDate(0, 0, day_index).Format("Jan 2 06")
	}

	file, err := os.Open("player.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	player_data, _ := reader.ReadAll()
	file.Close()

	var players []string
	for _, line := range player_data {
		players = append(players, line[0])
	}

	file, err = os.Open("wordle.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader = csv.NewReader(file)
	wordle_records, _ := reader.ReadAll()
	file.Close()

	var wordle_data []WordleData
	for _, line := range wordle_records {
		day, _ := strconv.Atoi(line[4])
		score, _ := strconv.Atoi(line[5])
		data := WordleData{
			ID:             line[0],
			MessageID:      line[1],
			UserID:         line[2],
			ChannelID:      line[3],
			Day:            day,
			Score:          score,
			BlankCount:     line[6],
			YellowCount:    line[7],
			GreenCount:     line[8],
			FirstWordScore: line[9],
		}
		wordle_data = append(wordle_data, data)
	}

	// for each user
	earliest_score := int(days_since_wordle_began)
	var bar_data []BarData
	for _, player_id := range players {
		stats := make([]int, days_since_wordle_began)
		for stat_index := range stats {
			for _, wordle := range wordle_data {
				if stat_index == wordle.Day && wordle.UserID == player_id {
					stats[stat_index] = wordle.Score
					if stat_index < earliest_score {
						earliest_score = stat_index
					}
				}
			}
		}
		bar := BarData{
			Label:           player_id,
			BackgroundColor: fmt.Sprintf("#%s", randomHex(3)),
			Data:            stats,
		}
		bar_data = append(bar_data, bar)
	}

	for index, bars := range bar_data {
		bar_data[index].Data = bars.Data[earliest_score:]
	}

	newData := GraphData{
		Labels: calendar[earliest_score:],
		Bars:   bar_data,
	}
	c.JSON(http.StatusOK, newData)
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
	router.GET("/", helloWorld)
	router.GET("/graph", fakeData)
	router.Run(":8000")
}
