// main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	scanner "Portscan/Scanner"
)

// RequestBody defines the structure of the incoming POST request.
type RequestBody struct {
	Subdomain string `json:"subdomain"`
}

// ResponseBody defines the structure of the API response.
type ResponseBody struct {
	Subdomain string `json:"subdomain"`
	Results   string `json:"results"`
	Error     string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/scan", handleScan)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleScan processes the incoming POST request, runs the port scan, and returns the results.
func handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON request body
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if reqBody.Subdomain == "" {
		http.Error(w, "Subdomain is required", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var result string
	var scanErr error

	wg.Add(1)
	go func(subdomain string) {
		defer wg.Done()
		result, scanErr = scanner.Scan(subdomain)
	}(reqBody.Subdomain)

	wg.Wait()

	// Prepare the response
	response := ResponseBody{
		Subdomain: reqBody.Subdomain,
		Results:   result,
	}

	if scanErr != nil {
		response.Error = scanErr.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
