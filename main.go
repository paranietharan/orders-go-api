package main

import (
	"context"
	"fmt"
	"orders-api/application"
)

func main() {
	ctx := context.Background()
	app := application.NewApp()
	if err := app.Start(ctx); err != nil {
		fmt.Printf("Failed to start application: %v\n", err)
	}
}
