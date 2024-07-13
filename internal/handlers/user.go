package handlers

import (
	"authgo/internal/services"
	"authgo/internal/types"
	"authgo/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var avatarService = services.NewAvatarServices()
var userServices = services.NewUserServices()

func Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"msg": "Invalid request payload"})
		return
	}

	id := uuid.New()

	filename, err := avatarService.GenerateAvatar(req.FirstName, id.String())
	if err != nil {
		log.Printf("Error generating avatar: %v", err)
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	user, err := userServices.RegisterByEmail(
		id,
		req.Email,
		filename,
		req.FirstName,
		req.LastName,
		req.Password,
	)

	if err != nil {
		if errRemove := os.Remove(filename); errRemove != nil {
			log.Printf("Failed to remove avatar file: %v", errRemove)
		}
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}
