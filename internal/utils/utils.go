package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func EnsureDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create  directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if directiry exists: %w", err)
	}

	return nil
}
