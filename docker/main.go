package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// VerificationResponse is used to send the verification challenge back
type VerificationResponse struct {
	Verification string `json:"verification"`
}

// EventData represents the structure of the event data received in POST request
type EventData struct {
	Data struct {
		Events []struct {
			Target []struct {
				AlternateId string `json:"alternateId"`
			} `json:"target"`
			EventType string `json:"eventType"`
		} `json:"events"`
	} `json:"data"`
}

func main() {
	http.HandleFunc("/user-lifecycle", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r)
		case http.MethodPost:
			handlePostRequest(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	challenge := r.Header.Get("x-okta-verification-challenge")
	response := VerificationResponse{Verification: challenge}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Println("Event hook verification request received.")
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	// Define the expected authorization header value
	const expectedAuthHeader = "Basic MjE0YzZmNGMtenRxYi03ODU2LTl1cmktYzliZmFhNDg4Nzdm"

	// Retrieve the Authorization header from the request
	authHeader := r.Header.Get("Authorization")

	// Check if the Authorization header matches the expected value
	if authHeader != expectedAuthHeader {
		fmt.Println("Authorization header: ", authHeader)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var eventData EventData
	err := json.NewDecoder(r.Body).Decode(&eventData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := eventData.Data.Events[0].Target[0].AlternateId
	eventType := eventData.Data.Events[0].EventType
	log.Printf("Lifecycle event %s for user %s", eventType, user)

	w.WriteHeader(http.StatusOK)
}
