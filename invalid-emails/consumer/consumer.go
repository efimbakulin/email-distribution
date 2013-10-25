package main

import (
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/connection-string-builder"
	"github.com/streadway/amqp"
	"log"
)

type ErrorChannel chan uint64
type OkChannel chan uint64

type MessageHandler func([]byte) error
type ErrorHandler func(tag uint64, consumer *Consumer) error

type Consumer struct {
	config properties.Properties
	conn   *amqp.Connection
	ch     *amqp.Channel
	done   chan bool
}

func RequeueMessageOnError(tag uint64, consumer *Consumer) error {
	return consumer.ch.Nack(tag, false, true)
}

func SkipMessageOnError(tag uint64, consumer *Consumer) error {
	return nil
}

func NewConsumer(config properties.Properties) *Consumer {
	return &Consumer{
		config: config,
		done: make(chan bool),
	}
}

func (self *Consumer) Connect() error {
	connBuilder, err := connstring.CreateBuilder(connstring.ConnectionStringAmqp)
	connBuilder.Address(self.config.String("rabbitmq.addr", ""))
	connBuilder.Port(uint16(self.config.Int("rabbitmq.port", 5672)))
	connBuilder.Username(self.config.String("rabbitmq.username", ""))
	connBuilder.Password(self.config.String("rabbitmq.password", ""))

	self.conn, err = amqp.Dial(connBuilder.Build())
	if err != nil {
		return err
	}
	self.ch, err = self.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	if err = self.ch.QueueBind(
		self.config.String("rabbitmq.queue.name", ""),
		"",
		self.config.String("rabbitmq.exchange.name", ""),
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Faild to bind queue: %s", err)
	}

	return nil
}

func (self *Consumer) Serve(messageHandler MessageHandler, errorHandler ErrorHandler) error {
	deliveries, err := self.ch.Consume(
		self.config.String("rabbitmq.queue.name", ""), // name
		"",    // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	errChan := make(ErrorChannel)
	okChan := make(OkChannel)
	go func() {
		var tag uint64
		select {
		case tag = <-errChan:
			{
				errorHandler(tag, self)
			}
		case tag = <-okChan:
			{
				self.ch.Ack(tag, false)
			}
		}
	}()

	go func() {
		for d := range deliveries {
			if err := messageHandler(d.Body); err != nil {
				errChan <- d.DeliveryTag
				continue
			}
			okChan <- d.DeliveryTag
		}
		log.Printf("handle: deliveries channel closed")
		self.done <- true
	}()
	return nil
}

func (self *Consumer) Stop() error {
	if self.ch != nil {
		self.ch.Close()
	}
	if self.conn != nil {
		self.conn.Close()
	}
	<-self.done
	return nil
}
