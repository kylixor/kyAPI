package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type BarData struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	Data            []int  `json:"data"`
}

type GraphData struct {
	Labels []string  `json:"labels"`
	Bars   []BarData `json:"datasets"`
}

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func fakeData(c *gin.Context) {
	newData := GraphData{
		Labels: []string{"January", "February", "March"},
		Bars: []BarData{
			{
				Label:           "Wordle Score",
				BackgroundColor: "#336699",
				Data:            []int{10, 20, 30},
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
