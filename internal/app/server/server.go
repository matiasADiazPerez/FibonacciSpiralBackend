package server

import (
	"log"
	"net/http"

	"spiralmatrix/internal/app/auth"
	"spiralmatrix/internal/app/user"

	"github.com/go-chi/chi/v5"
)

func Start(userh *user.UserHandler, authHandler *auth.AuthHandler) {
	router := chi.NewRouter()
	initRoutes(router, userh, authHandler)
	log.Println("Serving in :8080")
	http.ListenAndServe(":8080", router)
}

func initRoutes(router *chi.Mux, userh *user.UserHandler, authHandler *auth.AuthHandler) {
	router.Post("/user", userh.HandleCreateUser)
	router.Post("/login", authHandler.LoginHandler)
	router.Route("/user", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware())
		r.Get("/", userh.HandleGetAllUsers)
		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", userh.HandleGetUser)
			r.Patch("/", userh.HandleChangePass)
			r.Delete("/", userh.HandleDeleteUser)
		})
	})
}
