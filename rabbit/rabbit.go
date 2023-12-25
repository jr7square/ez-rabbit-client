package rabbit

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitProducer struct {
	RChannel *amqp.Channel
	conn     *amqp.Connection
}

var Producer *RabbitProducer

func InitRabbit(dialUrl string) *RabbitProducer {
	conn, err := amqp.Dial(dialUrl)
	if err != nil {
		log.Fatalln("Failed to connect to rabbit", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln("Faile to get channel", err)
	}

	fmt.Println("Successfully connected to rabbit")
	return &RabbitProducer{
		RChannel: channel,
		conn:     conn,
	}
}

func (producer *RabbitProducer) Close() error {
	err := producer.RChannel.Close()
	if err != nil {
		return err
	}
	err = producer.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (producer *RabbitProducer) Send(exchange string, key string, message string) error {
	if err := producer.RChannel.QueueBind(key, key, exchange, false, nil); err != nil {
		return err
	}
	err := producer.RChannel.PublishWithContext(context.TODO(),
		exchange,
		key,
		false,
		false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		})
	return err
}
