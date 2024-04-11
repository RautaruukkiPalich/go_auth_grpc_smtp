package smtp

import (
	"log/slog"

	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/app/kafka"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/config"
)

type SMTPClient struct {
	log *slog.Logger

}

func New(log *slog.Logger, cfg *config.Config) *SMTPClient {
	return &SMTPClient{
		log: log,
	}
}

func (c *SMTPClient) Run(msgch chan kafka.Message) {
	c.log.Info("start SMTP client")

	for {
		select {
			case data, ok :=<- msgch:
				if !ok {
					continue
				}
				if err := c.SendMessage(&data); err != nil {
					c.log.Error("send message")
				}
			default:
				continue
		}
	}
	
}

func (c *SMTPClient) Stop() {
	panic("stop")
}

func (c *SMTPClient) SendMessage(msg *kafka.Message) error {
	const op = "app.smtp.app.SendMessage"
	log := c.log.With(slog.String("op", op))
	log.With(slog.String("email", msg.Payload.Email)).Info("send message")
	
	// TODO

	return nil
}