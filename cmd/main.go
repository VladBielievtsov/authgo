package main

import (
	"authgo/db"
	"authgo/internal/config"
	"authgo/internal/handlers"
	"authgo/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		utils.ErrorHandler(err)
	}

	if err := db.ConnectDatabase(); err != nil {
		utils.ErrorHandler(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/register", handlers.Register)

	http.ListenAndServe(":"+cfg.Application.Port, r)
}
