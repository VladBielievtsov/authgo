package handlers

import (
	"authgo/internal/types"
	"authgo/internal/utils"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

}
