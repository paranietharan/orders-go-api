package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderId     uint64     `json:"orderId"`
	CustomerID  uuid.UUID  `json:"customerId"`
	LineItems   []LineItem `json:"lineItems"`
	CreatedAt   *time.Time `json:"createdAt"`
	ShippedAt   *time.Time `json:"shippedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"itemId"`
	Price    float64   `json:"price"`
	Quantity float64   `json:"quantity"`
}
