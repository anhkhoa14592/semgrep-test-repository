package main

// View layer. A Kafka consumer has no HTTP response to render, so "View"
// here means: the observable outcome of processing a message (logging).

import (
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func RenderProcessed(msg kafka.Message) {
	log.Printf("processed message: topic=%s partition=%d offset=%d", msg.Topic, msg.Partition, msg.Offset)
}

func RenderProcessingError(msg kafka.Message, reason interface{}) {
	log.Printf("failed to process message: topic=%s reason=%v", msg.Topic, reason)
}

func RenderUnroutedTopic(msg kafka.Message) {
	log.Printf("no route registered for topic=%s", msg.Topic)
}
