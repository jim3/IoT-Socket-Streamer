package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Define a struct to hold the sensor values
type SensorData struct {
	Temperature string
	Humidity    string
	Pressure    string
	Altitude    string
}

// Global instance of the struct and a mutex for thread-safe access
var sensorData SensorData
var mu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ------------------------------------------------------------

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		mu.Lock()
		err = conn.WriteJSON(sensorData)
		mu.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// ------------------------------------------------------------

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. Base template must be the *first* file in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Parse the files in the slice and store the resulting templates in a template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error in ParseFiles", 500)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base" template as the response body
	mu.Lock()
	defer mu.Unlock()
	err = ts.ExecuteTemplate(w, "base", sensorData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error in ExecuteTemplate", 500)
	}
}

// ------------------------------------------------------------

func createSensorData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest) // Response status: 400 Bad Request
		return
	}
	defer r.Body.Close()

	// Convert the body to a string
	bodyString := string(body)
	log.Println("Received data:", bodyString)

	// Split the string by spaces to separate the sensor values
	sensorValue := strings.Fields(bodyString)
	if len(sensorValue) < 4 {
		http.Error(w, "Invalid data format", http.StatusBadRequest) // Response body: Invalid data format
		return
	}

	// Extract the sensor values
	mu.Lock()
	sensorData = SensorData{
		Temperature: sensorValue[0],
		Humidity:    sensorValue[1],
		Pressure:    sensorValue[2],
		Altitude:    sensorValue[3],
	}
	mu.Unlock()

	log.Println("Sensor values:", sensorData)

	// Send a response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sensor data received successfully"))
}

// ------------------------------------------------------------

func main() {
	// Create a new ServeMux and register the routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc("/api/v1/sensors", createSensorData)

	// Start the server
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
