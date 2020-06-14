package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	numbers, err := readInput("input.txt")
	failOnError(err, "Failed to read input from file")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	dataQueue, err := ch.QueueDeclare(
		"data", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare data queue")

	body, err := json.Marshal(numbers)
	failOnError(err, "Failed to convert array of number to json message")

	err = ch.Publish(
		"",     // exchange
		dataQueue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	resultQueue, err := ch.QueueDeclare(
		"result", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare data queue")

	msgs, err := ch.Consume(
		resultQueue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go consume(msgs)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func consume(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		var fibonacciNumbers []int64
		err := json.Unmarshal(d.Body, &fibonacciNumbers)
		failOnError(err, "Failed to decode message")

		err = writeOutput("output.txt", fibonacciNumbers)
		failOnError(err, "Failed to write output to file")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
