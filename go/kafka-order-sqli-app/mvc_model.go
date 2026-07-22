package main

// Model layer.

import (
	"database/sql"
	"fmt"
)

type OrderModel struct {
	DB *sql.DB
}

func NewOrderModel(db *sql.DB) *OrderModel {
	return &OrderModel{DB: db}
}

// UpdateOrderStatus builds a raw SQL statement by interpolating orderID
// and status directly into the query text via fmt.Sprintf. Both values
// come straight from a Kafka message (see mvc_controller.go) - trusting
// the queue's payload implicitly is exactly the mistake this sample
// exists to flag. If a producer is compromised, or simply forwards
// unvalidated end-user input onto the topic, this is SQL injection.
func (m *OrderModel) UpdateOrderStatus(orderID string, status string) {
	query := fmt.Sprintf("UPDATE orders SET status = '%s' WHERE order_id = '%s'", status, orderID)
	m.DB.Exec(query)
}

// UpdateOrderStatusSafe is the safe counterpart - bind parameters instead
// of string interpolation. Kept as a precision check: the taint rule must
// NOT flag this path.
func (m *OrderModel) UpdateOrderStatusSafe(orderID string, status string) {
	m.DB.Exec("UPDATE orders SET status = ? WHERE order_id = ?", status, orderID)
}
