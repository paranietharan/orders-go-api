package application

import (
	"orders-api/handler"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/order", loadOrderRoutes)

	return router
}

func loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.OrderHandler{}

	r.Post("/", orderHandler.Create)
	r.Get("/", orderHandler.List)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.Update)
}
