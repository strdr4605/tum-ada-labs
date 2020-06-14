package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func main() {
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
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		dataQueue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	resultQueue, err := ch.QueueDeclare(
		"result", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare data queue")

	forever := make(chan bool)
	result := make(chan []int64)

	go consume(result, msgs)
	go sendResult(result, ch, resultQueue)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func consume(result chan []int64, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		var numbers []int
		err := json.Unmarshal(d.Body, &numbers)
		failOnError(err, "Failed to decode message")

		fibonacciNumbers := calculateFibonacci(numbers)
		result <- fibonacciNumbers
	}
}

func calculateFibonacci(numbers []int) []int64 {
	var fibonacciNumbers []int64
	for _, value := range numbers {
		fibonacciNumbers = append(fibonacciNumbers, sleepyFibonacci(value))
	}

	return fibonacciNumbers
}

func sendResult(resultChannel chan []int64, ch *amqp.Channel, q amqp.Queue) {
	for {
		select {
		case result := <- resultChannel:
			body, err := json.Marshal(result)
			failOnError(err, "Failed to convert array of number to json message")
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing {
					ContentType: "text/plain",
					Body:        body,
				})
			failOnError(err, "Failed to publish a message")
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
