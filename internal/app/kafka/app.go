package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"

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
	Message string `json:"message"`
}

type Broker struct {
	broker *kafka.Consumer
	log    *slog.Logger
	done   chan struct{}
}

type Consumer interface {
	Run() chan Message
	Stop()
}

func New(log *slog.Logger, cfg *config.Config) *Broker {

	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":  fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port),
			"group.id":           "kafka-consumer",
			"auto.offset.reset":  "smallest",
			"enable.auto.commit": "true",
		},
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(cfg.Kafka.Topic, nil)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	return &Broker{
		broker: c,
		log:    log,
		done:   done,
	}
}

func (b *Broker) Run() chan Message {
	const op = "app.kafka.app.GetFromQueue"
	log := b.log.With(slog.String("op", op))

	msgch := make(chan Message)

	go func() {
		for {
			ev := b.broker.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				var data Message
				err := json.Unmarshal(e.Value, &data)
				if err != nil{
					log.Error("error unmarshalling", op, slerr.Err(err))
				}
				msgch <- data
				// fmt.Printf("<<= Message: %s\n", string(e.Value))
			case kafka.Error:
				log.Error("error while sending message", slerr.Err(e))
			}
		}
	}()

	go func() {
		<-b.done
		close(msgch)
	}()

	return msgch
}

func (b *Broker) Stop() {
	close(b.done)
	b.broker.Close()
}
