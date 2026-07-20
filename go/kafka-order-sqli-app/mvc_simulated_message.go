package main

// Router-less variant: instead of a live broker connection plus
// TopicRouter, this builds a kafka.Message exactly as it would arrive on
// the wire (Key/Value only) and dispatches it straight to the controller
// - no router, no consume loop, no network needed. Useful both as the
// "without router" MVC variant and as a way to exercise the vulnerable
// path without standing up a real Kafka broker.

import (
	"database/sql"

	kafka "github.com/segmentio/kafka-go"
)

func SimulateIncomingOrderEvent(db *sql.DB, orderID string, status string) {
	msg := kafka.Message{
		Topic: "order-events",
		Key:   []byte(orderID),
		Value: []byte(status),
	}

	controller := NewOrderEventController(db)
	controller(msg)
}
