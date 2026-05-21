package database

import (
	"log"

	"his/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg *config.Config) *amqp.Connection {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v", err)
	}

	log.Println("Connected to RabbitMQ")

	return conn
}
