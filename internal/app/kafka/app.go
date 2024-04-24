package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/config"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/lib/slerr"
)

type Message struct {
	Topic   string  `json:"topic"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Email   string `json:"email"`
	Header  string `json:"header"`
	Message string `json:"message"`
}

type Broker struct {
	broker *kafka.Consumer
	log    *slog.Logger
	done   chan struct{}
	ticker *time.Ticker
}

type Consumer interface {
	Run() chan Payload
	Stop()
}

const (
	timeout = time.Millisecond * 1000
)

func New(log *slog.Logger, cfg *config.KafkaConfig) *Broker {
	const op = "app.kafka.app.Run"
	log.With(slog.String("op", op)).Info("start kafka consumer")

	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":  fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			"group.id":           cfg.ConsumerGroup,
			"auto.offset.reset":  "smallest",
			"enable.auto.commit": "true",
		},
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(cfg.Topic, nil)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(timeout)
	done := make(chan struct{})

	return &Broker{
		broker: c,
		log:    log,
		done:   done,
		ticker: ticker,
	}
}

func (b *Broker) Run() chan Payload {
	const op = "app.kafka.app.Run"
	log := b.log.With(slog.String("op", op))

	msgch := make(chan Payload)

	go func() {
		for {
			select {
			case <- b.done:
				return
			case <- b.ticker.C:
				ev := b.broker.Poll(100)
				switch e := ev.(type) {
				case *kafka.Message:
					var data Message
					err := json.Unmarshal(e.Value, &data)
					if err != nil {
						log.Error("error unmarshalling", op, slerr.Err(err))
					}
					msgch <- data.Payload
				case kafka.Error:
					log.Error("error while sending message", slerr.Err(e))
				}
			}
		}
	}()

	return msgch
}

func (b *Broker) Stop() {
	const op = "app.kafka.app.Stop"
	log := b.log.With(slog.String("op", op))
	log.Info("stop kafka client")

	close(b.done)
	defer b.ticker.Stop()
	defer b.broker.Close()
}
