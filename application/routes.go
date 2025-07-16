package application

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"orders-api/handler"
	"orders-api/repository/order"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/order", a.loadOrderRoutes)

	a.router = router
}

func (a *App) loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.OrderHandler{
		Repo: &order.RedisRepository{
			Redis: a.rdb,
		},
	}

	r.Post("/", orderHandler.Create)
	r.Get("/", orderHandler.List)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.Update)
}
