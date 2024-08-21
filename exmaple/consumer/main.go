package consumer

// create a main function and use confluent-kafka-go to subscribe to the topic and consume messages
import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"opentelemetry.io/contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/v2/otelkafka"
	"os"
)

func main() {
	// create a new consumer instance
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Created Consumer")

	otelConsumer := otelkafka.NewConsumer(consumer)

	// subscribe to the topic
	err = otelConsumer.SubscribeTopics([]string{"myTopic"}, nil)
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
			ev := otelConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}
