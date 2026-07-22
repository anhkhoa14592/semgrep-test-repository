package main

// Controller layer.
//
// NewOrderEventController is where a message "from Kafka" first touches
// this codebase - the consumer-side equivalent of an HTTP controller.
// orderID and status are read straight off the Kafka message (Key/Value)
// and handed to the Model layer unmodified: if either is producer/
// attacker controlled, this is the taint source for the SQL injection in
// mvc_model.go (OrderModel.UpdateOrderStatus).

import (
	"database/sql"

	kafka "github.com/segmentio/kafka-go"
)

func NewOrderEventController(db *sql.DB) func(msg kafka.Message) {
	model := NewOrderModel(db)

	return func(msg kafka.Message) {
		// SOURCE: values read directly off the Kafka message.
		orderID := string(msg.Key)
		status := string(msg.Value)

		if orderID == "" || status == "" {
			RenderProcessingError(msg, "missing order id or status")
			return
		}

		model.UpdateOrderStatus(orderID, status)
		RenderProcessed(msg)
	}
}

// NewOrderEventSafeController is the contrasting, non-vulnerable path
// through the same controller/model layers - kept so the taint rule's
// precision can be checked (it must NOT flag this one).
func NewOrderEventSafeController(db *sql.DB) func(msg kafka.Message) {
	model := NewOrderModel(db)

	return func(msg kafka.Message) {
		orderID := string(msg.Key)
		status := string(msg.Value)

		if orderID == "" || status == "" {
			RenderProcessingError(msg, "missing order id or status")
			return
		}

		model.UpdateOrderStatusSafe(orderID, status)
		RenderProcessed(msg)
	}
}
