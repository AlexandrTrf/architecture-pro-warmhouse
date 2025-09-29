package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

// Обработка /temperature/{sensorID}
func temperatureBySensorHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "sensorID missing in path", http.StatusBadRequest)
		return
	}
	sensorID := pathParts[2]
	location := ""

switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}

	resp := generateTemperature(sensorID, location)
	writeJSON(w, resp)
}

// Обработка /temperature?location=...
func temperatureByLocationHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "location query missing", http.StatusBadRequest)
		return
	}

    sensorID :=""
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}

	resp := generateTemperature(sensorID, location)
	writeJSON(w, resp)
}

// Генерация случайной температуры
func generateTemperature(sensorID, location string) TemperatureResponse {
	rand.Seed(time.Now().UnixNano())
	value := 15.0 + rand.Float64()*(30.0-15.0)

	return TemperatureResponse{
		Value:       value,
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Simulated temperature sensor",
	}
}

func writeJSON(w http.ResponseWriter, resp TemperatureResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/temperature/", temperatureBySensorHandler)
	http.HandleFunc("/temperature", temperatureByLocationHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
