package publisher

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"rabbitmq/util"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) *Publisher {
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	publisher := &Publisher{
		ch: ch,
	}
	go func() {
		publisher.Run()
	}()

	return publisher
}

func (p *Publisher) Run() {
	err := p.ch.ExchangeDeclare(
		"hieu_hoc_code", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
}

func (p *Publisher) Emit(c *gin.Context) {

	msg := "message"
	err := p.ch.Publish(
		"hieu_hoc_code",
		"routing_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	util.FailOnError(err, "Failed to publish a message")
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
