package messaging

import (
	"context"
	"encoding/json"
	"time"

	"his/internal/dto"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	AuditExchange = "his.audit.exchange"
	AuditQueue    = "his.audit.queue"
	AuditRouteKey = "audit.log"
)

type AuditPublisher interface {
	PublishAuditLog(ctx context.Context, event dto.AuditLogEvent) error
}

type RabbitMQAuditPublisher struct {
	conn *amqp.Connection
}

func NewRabbitMQAuditPublisher(conn *amqp.Connection) *RabbitMQAuditPublisher {
	return &RabbitMQAuditPublisher{
		conn: conn,
	}
}

func (p *RabbitMQAuditPublisher) PublishAuditLog(ctx context.Context, event dto.AuditLogEvent) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		AuditExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		AuditQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.QueueBind(
		AuditQueue,
		AuditRouteKey,
		AuditExchange,
		false,
		nil,
	); err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		AuditExchange,
		AuditRouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
