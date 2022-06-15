package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Hello World 2.0!")
}

func fakeData(w http.ResponseWriter, r *http.Request) {
	// var data GraphData
	// data.Labels = append(data.Labels, "January", "February", "March")
	// var barData BarData
	// barData.Label = "Wordle Score"
	// barData.BackgroundColor = "#f87979"
	// barData.Data = append(barData.Data, 40, 20, 12)
	// data.Bars = append(data.Bars, barData)
	newData := GraphData{
		Labels: []string{"January", "February", "March"},
		Bars: []BarData{
			{
				Label:           "Wordle Score",
				BackgroundColor: "#f87979",
				Data:            []int{40, 20, 12},
			},
		},
	}
	jsonResp, _ := json.Marshal(newData)
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/fake-data", fakeData)
	fmt.Println("Starting web server on 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
