package database

import (
	"log"
	"time"

	"his/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg *config.Config) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL)
		if err == nil {
			log.Println("Connected to RabbitMQ")
			return conn
		}

		log.Printf("RabbitMQ connection failed attempt %d/10: %v", i, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Unable to connect to RabbitMQ after retries: %v", err)

	return nil
}
