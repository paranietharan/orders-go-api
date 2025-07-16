package application

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func NewApp() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("could not connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Printf("Failed to close redis connection: %v\n", err)
		}
	}()

	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("Error starting server: %v", err)
		}
		close(ch)
	}()

	err = <-ch
	select {
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return server.Shutdown(timeoutCtx)
	case err = <-ch:
		return err
	}
	return nil
}
