package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	rabbitMQURL := "amqp://berna1:akman1@localhost:5672/"
	virtualHost := "to-do-api"
	queueName := "card_assignments"

	conn, err := amqp.Dial(fmt.Sprintf("%s%s", rabbitMQURL, virtualHost))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			card := Card{}
			err := json.Unmarshal(msg.Body, &card)
			if err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			err = sendEmail(card.Email, "Card Assigned", fmt.Sprintf("You have been assigned a new card:\nID: %s\nSummary: %s", card.ID, card.Summary))
			if err != nil {
				log.Printf("Failed to send email: %v", err)
				continue
			}

			log.Printf("Email sent to %s for card ID %s", card.Email, card.ID)
		}
	}()

	log.Println("Waiting for card assignment events. To exit, press CTRL+C")
	<-forever
}

func sendEmail(email, subject, body string) error {
	log.Printf("Sending email to %s\nSubject: %s\nBody: %s", email, subject, body)
	// Implement your email sending logic here
	return nil
}

type Card struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Email       string `json:"email"`
}
