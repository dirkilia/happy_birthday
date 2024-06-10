package server

import (
	"happy_birthday/internal/database"
	"happy_birthday/internal/middleware/jwtauth"
	"happy_birthday/internal/server/handlers/get"
	"happy_birthday/internal/server/handlers/post"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes(log *slog.Logger, db database.Service, jwtKey string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		// Only authorized users can access these handlers
		r.Use(jwtauth.CheckIfAuthorized)
		r.Get("/info", get.GetCurrentUserInfo(log, db))
		r.Post("/setnotify", post.SetNotifyOfHandler(log, db))
		r.Get("/emptdaybdays", get.GetAllEmployeesTodayBirthdaysHandler(log, db))
		r.Get("/employees", get.GetAllEmployeesHandler(log, db))
		r.Post("/logout", post.Logout(log, db))
	})

	r.Group(func(r chi.Router) {
		// Only unauthorized users can access these handlers
		r.Use(jwtauth.CheckIfNotAuthorized)
		r.Post("/register", post.RegisterUserHandler(log, db))
		r.Post("/login", post.AuthUserHandler(log, db))
	})

	return r
}
