package main

import (
	"encoding/json"
	"net/http"
)

type WeatherInfo struct {
	ID       int     `json:"id"`
	Temp     float64 `json:"temp"`
	Location string  `json:"location"`
	Humidity float64 `json:"humidity"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	weather := WeatherInfo{
		ID:       1,
		Temp:     22.5,
		Location: "San Francisco",
		Humidity: 60.0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	http.HandleFunc("/weather", weatherHandler)
	http.ListenAndServe(":8080", nil)
}
