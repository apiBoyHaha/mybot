package routes

import (
	"database/sql"
	"mybot/internal/auth"
	"mybot/internal/user"
	"net/http"
)

func SetupAuthRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	authHandler := auth.NewHandler(userService)

	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/refresh", authHandler.RefreshToken)

	return mux
}
