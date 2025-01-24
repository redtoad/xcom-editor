package main

import (
	"encoding/json"
	"net/http"
)

func locations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type point struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}
	type response []point

	resp := make(response, 0)
	for _, base := range currentSavegame.Bases() {
		coords := base.Coord()
		resp = append(resp, point{Latitude: coords.Lat, Longitude: coords.Lon})
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
