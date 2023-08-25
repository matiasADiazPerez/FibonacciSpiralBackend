package server

import (
	"log"
	"net/http"

	"spiralmatrix/internal/app/user"

	"github.com/go-chi/chi/v5"
)

func Start(userh *user.UserHandler) {
	router := chi.NewRouter()
	initRoutes(router, userh)
	log.Println("Serving in :8080")
	http.ListenAndServe(":8080", router)
}

func initRoutes(router *chi.Mux, userh *user.UserHandler) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/", userh.HandleGetAllUsers)
		r.Post("/", userh.HandleCreateUser)
		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", userh.HandleGetUser)
			r.Patch("/", userh.HandleChangePass)
			r.Delete("/", userh.HandleDeleteUser)
		})
	})
	router.Post("/login" auth.LoginHandler)
}
