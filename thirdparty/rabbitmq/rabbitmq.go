package rabbitmq

import (
	"github.com/streadway/amqp"
)

var consumer *RabbitMqHelper
var publisher *RabbitMqHelper

func InitConsumer(exchange, queue, addr, routingKey string) error {
	consumer = NewRabbitMqHelper(exchange, queue, addr, routingKey)
	return consumer.Connect()
}

func InitPublisher(exchange, queue, addr string) error {
	publisher = NewRabbitMqHelper(exchange, queue, addr, "")
	return publisher.Connect()
}

func Publish(routingKey string, msgContent []byte) error {
	return publisher.Publish(routingKey, msgContent)
}

func Consume() (<-chan amqp.Delivery, error) {
	return consumer.Consume()
}

type RabbitMQConfig struct {
	Exchange       string
	Queue          string
	Addr           string
	BindRoutingKey string //作为消费者需要的字段
}

type RabbitMqHelper struct {
	cf      *RabbitMQConfig
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMqHelper(exchange, queue, addr, routingKey string) *RabbitMqHelper {
	return &RabbitMqHelper{
		cf: &RabbitMQConfig{
			Exchange:       exchange,
			Queue:          queue,
			Addr:           addr,
			BindRoutingKey: routingKey,
		},
	}
}

func (this *RabbitMqHelper) Reconnect() (err error) {
	if this.conn == nil || this.conn.IsClosed() {
		this.conn, err = amqp.Dial(this.cf.Addr)
		if err != nil {
			return err
		}

		this.channel, err = this.conn.Channel()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *RabbitMqHelper) Close() {
	_ = this.conn.Close()
	_ = this.channel.Close()
}

func (this *RabbitMqHelper) Connect() error {
	var err error
	err = this.Reconnect()
	if err != nil {
		return err
	}

	err = this.declareExchangeAndQueue(amqp.ExchangeTopic, "")
	if err != nil {
		return err
	}
	return err
}

func (this *RabbitMqHelper) Publish(routingKey string, b []byte) error {
	return this.publish(routingKey, b)
}

func (this *RabbitMqHelper) Consume() (<-chan amqp.Delivery, error) {
	var err error
	err = this.Reconnect()
	if err != nil {
		return nil, err
	}

	err = this.declareExchangeAndQueue(amqp.ExchangeTopic, this.cf.BindRoutingKey)
	if err != nil {
		return nil, err
	}

	msgs, err := this.channel.Consume(
		this.cf.Queue, // fanoutQueue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	return msgs, err
}

func (this *RabbitMqHelper) declareExchangeAndQueue(mode string, bindRouteKey string) error {
	err := this.channel.ExchangeDeclare(
		this.cf.Exchange,
		mode,
		true,
		false,
		false,
		false,
		amqp.Table{})
	if err != nil {
		return err
	}

	q, err := this.channel.QueueDeclare(
		this.cf.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = this.channel.QueueBind(q.Name, bindRouteKey, this.cf.Exchange, false, amqp.Table{})
	return err
}

func (this *RabbitMqHelper) publish(routingKey string, msgContent []byte) error {
	err := this.channel.Publish(
		this.cf.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         msgContent,
			DeliveryMode: amqp.Persistent,
		})

	return err
}
