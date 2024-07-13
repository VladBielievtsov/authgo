package handlers

import (
	"authgo/internal/services"
	"authgo/internal/types"
	"authgo/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := uuid.New()

	filename, err := services.GenerateAvatar(req.FirstName, id.String())
	if err != nil {
		log.Printf("Error generating avatar: %v", err)
		utils.JSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"mgs": "Registration successful", "avatar": filename})
}
