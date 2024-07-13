package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
