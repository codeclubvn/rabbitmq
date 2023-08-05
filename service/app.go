package service

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"rabbitmq/conf"
	"rabbitmq/consumer"
	"rabbitmq/publisher"
)

type App struct {
	Router    *gin.Engine
	publisher *publisher.Publisher
	consumer  *consumer.Consumer
	conn      *amqp.Connection
}

func NewApp() *App {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()

	routesConfig, err := conf.GetConfigRoutingKey()
	if err != nil {
		panic(err)
	}

	return &App{
		Router:    gin.Default(),
		publisher: publisher.NewPublisher(conn),
		conn:      conn,
		consumer:  consumer.NewConsumer(conn, routesConfig),
	}
}

func (a *App) Start() error {
	return a.Router.Run(":" + conf.GetEnv().Port)
}

func (a *App) Emit(c *gin.Context) {
	a.publisher.Emit(c)
}
