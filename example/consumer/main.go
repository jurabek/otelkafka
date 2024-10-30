package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"github.com/jurabek/otelkafka/example"
	"go.opentelemetry.io/otel"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	topic := os.Getenv("KAFKA_TOPIC")
	kafkaServers := os.Getenv("KAFKA_SERVER")

	tp, err := example.InitTracer("consumer-app")
	if err != nil {
		log.Fatal(err)
	}

	mp, err := example.InitMeter("consumer-app")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}

		_ = mp.Shutdown(context.Background())
	}()

	consumer, err := otelkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":      kafkaServers,
		"group.id":               "myGroup",
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 5000,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	// subscribe to the topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Subscribed to myTopic")

	// consume messages
	run := true
	for run == true {
		select {
		case sig := <-signals:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				parentSpanContext := otel.GetTextMapPropagator().Extract(context.Background(), otelkafka.NewMessageCarrier(e))
				fmt.Printf("span context: %v\n", parentSpanContext)

			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}
