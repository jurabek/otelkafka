package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"github.com/jurabek/otelkafka/example"
	"go.opentelemetry.io/otel"
	"log"
	"os"
	"time"
)

func main() {
	tp, err := example.InitTracer("producer-app")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	bootstrapServers := os.Getenv("KAFKA_SERVER")
	topic := os.Getenv("KAFKA_TOPIC")

	p, err := otelkafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Producer %v\n", p)

	tr := otel.Tracer("produce")
	ctx, span := tr.Start(context.Background(), "produce message")
	defer span.End()

	deliveryChan := make(chan kafka.Event)

	for i := 0; i < 500; i++ {
		// Optional delivery channel, if not specified the Producer object's
		// .Events channel is used.
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(fmt.Sprintf("Hello Go #%d", i)),
			Key:            []byte("message-key"),
			Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}
		otel.GetTextMapPropagator().Inject(ctx, otelkafka.NewMessageCarrier(message))

		err = p.Produce(message, deliveryChan)
		if err != nil {
			fmt.Printf("Produce failed: %v\n", err)
			os.Exit(1)
		}

		e := <-deliveryChan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
		time.Sleep(5 * time.Second)
	}

	close(deliveryChan)
}
