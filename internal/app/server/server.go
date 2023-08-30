package server

import (
	"log"
	"net/http"

	"spiralmatrix/internal/app/auth"
	"spiralmatrix/internal/app/spiral"
	"spiralmatrix/internal/app/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Start(userh *user.UserHandler, authHandler *auth.AuthHandler, spiralHandler *spiral.SpiralHandler) {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	initRoutes(router, userh, authHandler, spiralHandler)
	log.Println("Start server in port :8080")
	http.ListenAndServe(":8080", router)
}

func initRoutes(router *chi.Mux, userh *user.UserHandler, authHandler *auth.AuthHandler, spiralHandler *spiral.SpiralHandler) {
	router.Post("/public/user", userh.HandleCreateUser)
	router.Post("/login", authHandler.LoginHandler)
	router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/doc.json"),
	))
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
