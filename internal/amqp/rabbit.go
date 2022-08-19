package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YuriyNazarov/bannersRotator/internal/config"
	"time"

	goamqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	exchange string
	queue    string
	consumer string
	channel  *goamqp.Channel
	logger   Logger
}

func (q *Rabbit) Click(bannerId, slotId, groupId int, clickTime time.Time) {
	msg := statsMessage{
		BannerId:   bannerId,
		SlotId:     slotId,
		GroupId:    groupId,
		Timestamp:  clickTime,
		ActionType: "click",
	}
	err := q.add(context.Background(), msg)
	if err != nil {
		q.logger.Error(fmt.Sprintf("err on sending click event: %s", err))
	}
}

func (q *Rabbit) Show(bannerId, slotId, groupId int, clickTime time.Time) {
	msg := statsMessage{
		BannerId:   bannerId,
		SlotId:     slotId,
		GroupId:    groupId,
		Timestamp:  clickTime,
		ActionType: "show",
	}
	err := q.add(context.Background(), msg)
	if err != nil {
		q.logger.Error(fmt.Sprintf("err on sending show event: %s", err))
	}
}

func NewRabbit(ctx context.Context, logger Logger, cfg config.QueueCfg) *Rabbit {
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

func (q *Rabbit) add(ctx context.Context, msg statsMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal banner event: %w", err)
	}

	err = q.channel.PublishWithContext(
		ctx,
		q.exchange,
		q.queue,
		false,
		false,
		goamqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		})
	if err != nil {
		return fmt.Errorf("failed send statistics: %w", err)
	}

	return nil
}
