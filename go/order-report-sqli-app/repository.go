package main

// Model / repository layer — chứa sink.

import (
	"database/sql"
	"fmt"
)

type Order struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) ListByStatus(status string) ([]Order, error) {
	return queryOrdersByStatus(r.db, status)
}

// TRUE POSITIVE: status đi thẳng từ query param vào chuỗi SQL qua fmt.Sprintf.
func queryOrdersByStatus(db *sql.DB, status string) ([]Order, error) {
	query := fmt.Sprintf("SELECT id, status, total FROM orders WHERE status = '%s'", status)
	rows, err := db.Query(query) // SINK
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Status, &o.Total); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
