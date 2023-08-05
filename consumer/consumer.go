package consumer

import (
	"bytes"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"rabbitmq/conf"
	"rabbitmq/util"
)

type Consumer struct {
	ch           *amqp.Channel
	routesConfig conf.RoutesConfig
	routingKey   map[string]string
}

func NewConsumer(conn *amqp.Connection, routesConfig conf.RoutesConfig) *Consumer {
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	consumer := &Consumer{
		ch:           ch,
		routesConfig: routesConfig,
		routingKey:   make(map[string]string),
	}
	go func() {
		consumer.Run()
	}()

	return consumer
}

func (c *Consumer) Run() {
	err := c.ch.ExchangeDeclare(
		"hieu_hoc_code", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	q, err := c.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	// 2 vòng for này có thể tối ưu
	for _, s := range c.routesConfig.Routes {
		c.routingKey[s.RoutingKey] = s.API

		err = c.ch.QueueBind(
			q.Name,          // queue name
			s.RoutingKey,    // routing key
			"hieu_hoc_code", // exchange
			false,
			nil)
		util.FailOnError(err, "Failed to bind a queue")
	}

	msgs, err := c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			bodyReader := bytes.NewReader(d.Body)
			_, err = http.Post(conf.GetEnv().ERPDomain+c.routingKey[d.RoutingKey], conf.GetEnv().ContentType, bodyReader)
			if err != nil {
				fmt.Println("cannot call to " + conf.GetEnv().ERPDomain)
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
