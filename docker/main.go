package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// VerificationResponse is used to send the verification challenge back
type VerificationResponse struct {
	Verification string `json:"verification"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Log the raw request body for debugging
	log.Println("Raw request body:", string(bodyBytes))

	// Re-create the request body for JSON decoding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var eventData EventData
	err = json.NewDecoder(r.Body).Decode(&eventData)
	if err != nil {
		log.Println("Error decoding JSON body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, event := range eventData.Data.Events {
		eventType := event.EventType
		switch eventType {
		case GroupLifecycleCreate:
			GroupLifecycleCreateHandler(event)
		case GroupMembershipAdd:
			GroupMembershipAddHandler(event)
		case GroupMembershipRemove:
			GroupMembershipRemoveHandler(event)
		case GroupProfileUpdate:
			GroupProfileUpdateHandler(event)
		case GroupApplicationAssignmentAdd:
			GroupApplicationAssignmentAddHandler(event)
		case GroupApplicationAssignmentRemove:
			GroupApplicationAssignmentRemoveHandler(event)
		default:
			log.Println("Unknown event type received.")
		}
	}
	w.WriteHeader(http.StatusOK)
}
