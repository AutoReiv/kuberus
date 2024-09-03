package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON serializes the given data into JSON and writes it to the HTTP response.
func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
