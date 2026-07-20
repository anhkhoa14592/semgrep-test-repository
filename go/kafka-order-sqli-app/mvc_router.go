package main

// Router layer (optional - see mvc_simulated_message.go for the
// router-less variant of the same flow).
//
// TopicRouter is the Kafka analogue of an HTTP router: instead of mapping
// a path to a handler, it maps a topic name to the controller responsible
// for it.

import (
	"context"
	"database/sql"

	kafka "github.com/segmentio/kafka-go"
)

type TopicRouter struct {
	routes map[string]func(kafka.Message)
}

func NewRouter(db *sql.DB) *TopicRouter {
	return &TopicRouter{
		routes: map[string]func(kafka.Message){
			"order-events": NewOrderEventController(db),
		},
	}
}

func (rt *TopicRouter) Dispatch(msg kafka.Message) {
	if h, ok := rt.routes[msg.Topic]; ok {
		h(msg)
		return
	}
	RenderUnroutedTopic(msg)
}

// RunWithRouter is the "message from Kafka" equivalent of ListenAndServe:
// it continuously reads messages off the broker and hands each one to the
// router. Called from main.go.
func RunWithRouter(ctx context.Context, reader *kafka.Reader, db *sql.DB) error {
	router := NewRouter(db)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		router.Dispatch(msg)
	}
}
