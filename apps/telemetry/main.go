package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Sensor struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	Unit        string    `json:"unit"`
	Value       float64   `json:"value"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`
}

type Telemetry struct {
	ID         string                 `json:"id"`
	SensorID   int                    `json:"sensor_id"`
	SensorType string                 `json:"sensor_type"`
	Value      float64                `json:"value"`
	Unit       string                 `json:"unit"`
	Location   string                 `json:"location"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type Event struct {
	EventID   string                 `json:"event_id"`
	EventType string                 `json:"event_type"`
	SensorID  int                    `json:"sensor_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Severity  string                 `json:"severity"`
}

var (
	sensors          = make(map[int]*Sensor)
	telemetry        = []Telemetry{}
	events           = []Event{}
	mu               sync.RWMutex
	temperatureAPIURL = getEnv("TEMPERATURE_API_URL", "http://temperature-api:8080")
	rabbitURL        = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
)

// Helper to get env variable or default
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// --- RabbitMQ Helper ---
func publishSensorEvent(s *Sensor) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sensor_events", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"event_type": "sensor_created",
		"sensor":     s,
		"timestamp":  time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Printf("Published sensor_created event for sensor %d", s.ID)
	}
}

// --- Handlers ---

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getSensorsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	list := []*Sensor{}
	for _, s := range sensors {
		list = append(list, s)
	}
	json.NewEncoder(w).Encode(list)
}

func createSensorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Location string  `json:"location"`
		Unit     string  `json:"unit"`
		Value    float64 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	id := len(sensors) + 1
	if input.Unit == "" {
		input.Unit = "Celsius"
	}
	s := &Sensor{
		ID:          id,
		Name:        input.Name,
		Type:        input.Type,
		Location:    input.Location,
		Unit:        input.Unit,
		Value:       input.Value,
		Status:      "active",
		LastUpdated: time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}
	sensors[id] = s

	// Publish event to RabbitMQ
	go publishSensorEvent(s)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func getSensorHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/sensors/"):]
	id, _ := strconv.Atoi(idStr)
	mu.RLock()
	defer mu.RUnlock()
	if s, ok := sensors[id]; ok {
		json.NewEncoder(w).Encode(s)
		return
	}
	http.Error(w, "Sensor not found", http.StatusNotFound)
}

func deleteSensorHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/sensors/"):]
	id, _ := strconv.Atoi(idStr)
	mu.Lock()
	defer mu.Unlock()
	if _, ok := sensors[id]; ok {
		delete(sensors, id)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sensor deleted successfully"})
		return
	}
	http.Error(w, "Sensor not found", http.StatusNotFound)
}

func updateSensorValueHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/sensors/"):]
	idStr = idStr[:len(idStr)-len("/value")]
	id, _ := strconv.Atoi(idStr)
	var input struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if s, ok := sensors[id]; ok {
		s.Value = input.Value
		s.Status = input.Status
		s.LastUpdated = time.Now().UTC()
		json.NewEncoder(w).Encode(map[string]string{"message": "Sensor value updated successfully"})
		return
	}
	http.Error(w, "Sensor not found", http.StatusNotFound)
}

// --- Background updater ---
func updateSensorValuesTask() {
	for {
		time.Sleep(1 * time.Minute)
		mu.Lock()
		for _, s := range sensors {
			url := temperatureAPIURL + "/temperature/" + strconv.Itoa(s.ID)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error updating sensor %d: %v", s.ID, err)
				continue
			}
			var data struct {
				Value  float64 `json:"value"`
				Status string  `json:"status"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				resp.Body.Close()
				log.Printf("Error parsing response for sensor %d: %v", s.ID, err)
				continue
			}
			resp.Body.Close()
			s.Value = data.Value
			s.Status = data.Status
			s.LastUpdated = time.Now().UTC()
		}
		mu.Unlock()
	}
}

// --- Main ---
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/sensors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getSensorsHandler(w, r)
		case "POST":
			createSensorHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/sensors/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			getSensorHandler(w, r)
		case r.Method == "DELETE":
			deleteSensorHandler(w, r)
		case r.Method == "PATCH" && len(r.URL.Path) > len("/api/v1/sensors/") && r.URL.Path[len(r.URL.Path)-6:] == "/value":
			updateSensorValueHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	go updateSensorValuesTask()

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", mux)
}
