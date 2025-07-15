package main

import (
	"context"
	"fmt"
	"orders-api/application"
	"os"
	"os/signal"
)

func main() {
	app := application.NewApp()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Printf("Failed to start application: %v\n", err)
	}
}
