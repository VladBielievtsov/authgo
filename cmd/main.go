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

	db.Migrate()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	fs := http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads")))
	r.Handle("/uploads/*", fs)

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	http.ListenAndServe(":"+cfg.Application.Port, r)
}
