package main

import (
	"encoding/csv"
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

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func fakeData(c *gin.Context) {

	days_since_wordle_began := uint16(time.Since(WORDLE_DAY_0).Hours() / 24)

	calendar := make([]string, days_since_wordle_began)
	for day_index := range calendar {
		calendar[day_index] = WORDLE_DAY_0.AddDate(0, 0, day_index).Format("Jan 2 06")
	}

	file, err := os.Open("wordle.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	wordle_data := []WordleData{}
	for _, line := range records {
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
	wordle_stats := make([]int, days_since_wordle_began)
	first_day := 0
	for stat_index := range wordle_stats {
		for _, wordle := range wordle_data {
			if stat_index == wordle.Day && wordle.UserID == "144220178853396480" {
				wordle_stats[stat_index] = wordle.Score
				if first_day == 0 {
					first_day = stat_index
				}
			}
		}
	}

	newData := GraphData{
		Labels: calendar[first_day:],
		Bars: []BarData{
			{
				Label:           "Andrew",
				BackgroundColor: "#336699",
				Data:            wordle_stats[first_day:],
			},
		},
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
	router.GET("/fake-data", fakeData)
	router.Run(":3000")
}
