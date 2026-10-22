package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSON writes payload as JSON with the given status.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if payload == nil {
		w.WriteHeader(status)
		return
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("api: marshal response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(body)
}

// WriteProblem writes problem as JSON with the given status.
func WriteProblem(w http.ResponseWriter, status int, problem any) {
	w.Header().Set("Content-Type", "application/json")
	if problem == nil {
		w.WriteHeader(status)
		return
	}
	body, err := json.Marshal(problem)
	if err != nil {
		log.Printf("api: marshal problem: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(body)
}
