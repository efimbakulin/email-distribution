package main

import (
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/fimston/connection-string-builder"
	"github.com/streadway/amqp"
)

type Producer struct {
	config             properties.Properties
	rabbitMqConnection *amqp.Connection
	rabbitMqChannel    *amqp.Channel
}

func (self *Producer) Connect() error {
	connBuilder, err := connstring.CreateBuilder(connstring.ConnectionStringAmqp)
	connBuilder.Address(self.config.String("rabbitmq.addr", ""))
	connBuilder.Port(uint16(self.config.Int("rabbitmq.port", 5672)))
	connBuilder.Username(self.config.String("rabbitmq.username", ""))
	connBuilder.Password(self.config.String("rabbitmq.password", ""))

	self.rabbitMqConnection, err = amqp.Dial(connBuilder.Build())
	if err != nil {
		return err
	}
	self.rabbitMqChannel, err = self.rabbitMqConnection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	return nil
}

func NewProducer(config properties.Properties) *Producer {
	instance := &Producer{}
	instance.config = config
	return instance
}

func (self *Producer) PostTask(data string) error {

	if err := self.rabbitMqChannel.Publish(
		self.config.String("rabbitmq.exchange.name", ""),
		self.config.String("rabbitmq.routing_key", ""),
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(data),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	); err != nil {
		return err
	}
	return nil
}

func (self *Producer) Stop() {
	if self.rabbitMqConnection != nil {
		self.rabbitMqConnection.Close()
	}
}
