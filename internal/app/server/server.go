package server

import (
	"log"
	"net/http"

	"spiralmatrix/internal/app/auth"
	"spiralmatrix/internal/app/spiral"
	"spiralmatrix/internal/app/user"

	"github.com/go-chi/chi/v5"
)

func Start(userh *user.UserHandler, authHandler *auth.AuthHandler, spiralHandler *spiral.SpiralHandler) {
	router := chi.NewRouter()
	initRoutes(router, userh, authHandler, spiralHandler)
	log.Println("Start server in port :8080")
	http.ListenAndServe(":8080", router)
}

func initRoutes(router *chi.Mux, userh *user.UserHandler, authHandler *auth.AuthHandler, spiralHandler *spiral.SpiralHandler) {
	router.Post("/public/user", userh.HandleCreateUser)
	router.Post("/login", authHandler.LoginHandler)
	router.With(authHandler.AuthMiddleware()).Group(func(r chi.Router) {
		r.Get("/spiral", spiralHandler.HandleSpiral)
		r.Route("/user", func(r chi.Router) {
			r.Get("/", userh.HandleGetAllUsers)
			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", userh.HandleGetUser)
				r.Patch("/", userh.HandleChangePass)
				r.Delete("/", userh.HandleDeleteUser)
			})
		})

	})
}
