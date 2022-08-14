package amqp

import (
	"context"
	"fmt"
	"github.com/YuriyNazarov/bannersRotator/internal/app"
	"github.com/YuriyNazarov/bannersRotator/internal/config"

	goamqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	exchange string
	queue    string
	consumer string
	channel  *goamqp.Channel
	logger   app.Logger
}

func NewRabbit(ctx context.Context, logger app.Logger, cfg config.QueueCfg) *Rabbit {
	conn, err := goamqp.Dial(cfg.DSN)
	if err != nil {
		logger.Error(fmt.Sprintf("failed on connect to rabblitmq: %s", err))
		return nil
	}

	chanel, err := conn.Channel()
	if err != nil {
		logger.Error(fmt.Sprintf("failed on opening chanel: %s", err))
		return nil
	}
	err = chanel.ExchangeDeclare(
		cfg.Exchange,
		goamqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("failed on creating exchange: %s", err))
		return nil
	}

	queue, err := chanel.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("failed on creating queue: %s", err))
		return nil
	}

	err = chanel.QueueBind(
		queue.Name,
		queue.Name,
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("failed on binding queue: %s", err))
		return nil
	}

	go func() {
		<-ctx.Done()
		chanel.Close()
		conn.Close()
	}()

	return &Rabbit{
		exchange: cfg.Exchange,
		queue:    cfg.Queue,
		consumer: "banners-stats-consumer",
		channel:  chanel,
		logger:   logger,
	}
}
