package main

import (
	"context"
	"encoding/json"
	"log"

	"his/internal/config"
	"his/internal/database"
	"his/internal/dto"
	"his/internal/messaging"
	"his/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewPostgres(cfg)
	rabbitConn := database.NewRabbitMQ(cfg)
	defer rabbitConn.Close()

	auditRepo := repository.NewAuditLogRepository(db)

	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		messaging.AuditExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	if _, err := ch.QueueDeclare(
		messaging.AuditQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	if err := ch.QueueBind(
		messaging.AuditQueue,
		messaging.AuditRouteKey,
		messaging.AuditExchange,
		false,
		nil,
	); err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		messaging.AuditQueue,
		"his-audit-worker",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Println("Audit worker started. Waiting for messages...")

	for msg := range msgs {
		var event dto.AuditLogEvent

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Println("Invalid audit event:", err)
			_ = msg.Nack(false, false)
			continue
		}

		if err := auditRepo.Create(context.Background(), event); err != nil {
			log.Println("Failed to save audit log:", err)
			_ = msg.Nack(false, true)
			continue
		}

		log.Println("Audit log saved:", event.EventType)
		_ = msg.Ack(false)
	}
}
