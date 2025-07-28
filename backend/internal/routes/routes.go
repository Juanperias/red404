package routes

import (
	"net/http"

	"github.com/escuadron-404/red404/backend/internal/handlers"
	"github.com/escuadron-404/red404/backend/pkg/middleware"
)

// SetupRoutes configures all application routes.
func SetupRoutes(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) http.Handler {
	mux := http.NewServeMux()

	// Register authentication-related routes
	AuthRoutes(mux, authHandler)

	// Register user-related routes
	UserRoutes(mux, userHandler, authMiddleware)

	return mux
}

func AuthRoutes(mux *http.ServeMux, authHandler *handlers.AuthHandler) {
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/register", authHandler.Register)
}

func UserRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// mux.HandleFunc("POST /api/users", authMiddleware.Auth(userHandler.CreateUser))
	mux.HandleFunc("GET /api/users/{id}", authMiddleware.Auth(userHandler.GetUserByID))
	mux.HandleFunc("GET /api/users", authMiddleware.Auth(userHandler.GetAllUsers))
	// mux.HandleFunc("PUT /api/users/{id}", authMiddleware.Auth(userHandler.UpdateUser))
	// mux.HandleFunc("DELETE /api/users/{id}", authMiddleware.Auth(userHandler.DeleteUser))
}
