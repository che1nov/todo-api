package main

import (
	"context"
	"log/slog"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/db"
	"todo-api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadConfig()

	dbInstance := db.InitDB()
	defer func() {
		sqlDB, _ := dbInstance.DB()
		sqlDB.Close()
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", dbInstance)
			ctx = context.WithValue(ctx, "config", cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	r.Group(func(r chi.Router) {
		r.Use(handlers.MiddlewareAuth(cfg.JWTSecret))
		r.Post("/todos", handlers.CreateTodo)
		r.Get("/todos", handlers.GetTodos)
		r.Get("/todos/{id}", handlers.GetTodo)
		r.Put("/todos/{id}", handlers.UpdateTodo)
		r.Delete("/todos/{id}", handlers.DeleteTodo)
	})

	slog.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
