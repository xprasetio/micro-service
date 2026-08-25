package message

import (
	"encoding/json"
	"user-service/config"

	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage(email, message, notif_type string) error {
	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		log.Errorf("[PublishMessage-1] failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[PublishMessage-2] failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		notif_type,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("[PublishMessage-3] failed to declare a queue: %v", err)
		return err
	}

	err = ch.QueueBind(
		queue.Name,
		notif_type,
		"amq.direct",
		false,
		nil,
	)
	if err != nil {
		log.Errorf("[PublishMessage-3] failed to bind a queue: %v", err)
		return err
	}

	notification := map[string]string{
		"email":   email,
		"message": message,
		"type":    notif_type,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		log.Errorf("[PublishMessage-3] failed to marshal a message: %v", err)
		return err
	}

	err = ch.Publish(
		"amq.direct", // exchange
		queue.Name,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		log.Errorf("[PublishMessage-3] failed to publish a message: %v", err)
		return err
	}
	return nil
}
