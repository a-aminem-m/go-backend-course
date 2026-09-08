package main

import (
	"encoding/json"
	"hard-hw1/models"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare(
		"code_tasks",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	for message := range messages {
		var task models.CodeTaskMessage
		err := json.Unmarshal(message.Body, &task)
		if err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}
		processTask(task)
	}
}
