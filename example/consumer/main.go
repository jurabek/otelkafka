package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jurabek/otelkafka"
	"github.com/jurabek/otelkafka/example"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"os"
	"os/signal"
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
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
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

				print(parentSpanContext, e)

			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case *kafka.Stats:
				// Stats events are emitted as JSON (as string).
				// Either directly forward the JSON to your
				// statistics collector, or convert it to a
				// map to extract fields of interest.
				// The definition of the statistics JSON
				// object can be found here:
				// https://github.com/confluentinc/librdkafka/blob/master/STATISTICS.md
				var stats map[string]interface{}
				json.Unmarshal([]byte(e.String()), &stats)
				// write stats into file
				fmt.Println(e.String())

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}

func print(ctx context.Context, msg *kafka.Message) {
	ctx, span := otel.Tracer("kafka_consumer").Start(ctx, "")
	span.SetAttributes(attribute.String("kafka.msg.key", string(msg.Key)))
	defer span.End()

	fmt.Printf("%% Message on %s:\n%s\n",
		msg.TopicPartition, string(msg.Key))
}
