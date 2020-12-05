package rabbitmq

import (
	"log"
	"testing"
	"time"
)

var (
	toppicExchange = "toppicExchange"
	toppicQueue    = "toppicQueue"
	rabbitMqAddr   = "amqp://admin:123456@8.131.67.245:5672/"
	bindRoutingKey = "a.*"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func TestRabbitMqHelperToppic(t *testing.T) {

	err := InitConsumer(toppicExchange, toppicQueue, rabbitMqAddr, bindRoutingKey)
	failOnError(err, "Failed to connect")

	err = InitPublisher(toppicExchange, toppicQueue, rabbitMqAddr)
	failOnError(err, "Failed to connect")

	err = Publish("a.1", []byte("Hello World from a.1!"))
	failOnError(err, "Failed to publish a message")

	err = Publish("a.2", []byte("Hello World from a.2!"))
	failOnError(err, "Failed to publish a message")

	err = Publish("b.1", []byte("Hello World from b.1!"))
	failOnError(err, "Failed to publish a message")

	msgs, err := Consume()
	failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	time.Sleep(1 * time.Second)
}
