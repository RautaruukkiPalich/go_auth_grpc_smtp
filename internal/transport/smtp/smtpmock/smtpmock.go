package smtpmock

import (
	"fmt"
	"log/slog"

	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/config"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/lib/slerr"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/model"
)

const (
	EmptyString = ""
)

var (
	ErrEmptyEmail = fmt.Errorf("empty email address")
	ErrEmptyMessage = fmt.Errorf("empty message")
)

type SMTPClient struct {
	log *slog.Logger
	done chan struct{}
	auth *slog.Logger	
	addr string
	from string
}

func New(log *slog.Logger, cfg *config.SMTPConfig) *SMTPClient {
	
	done := make(chan struct{})
	smtp := log.With(
		slog.Attr(
			slog.String("mock-client", "smtp"),
		),
	)

	return &SMTPClient{
		log: log,
		done: done,
		auth: smtp,
		addr: cfg.Addr,
		from: cfg.From,
	}
}

func (c *SMTPClient) Run(msgch chan model.Payload) {
	const op = "transport.smtp.mock.Run"
	log := c.log.With(slog.String("op", op))
	log.Info("start mock SMTP client")

	for {
		select {
			case data, ok :=<- msgch:
				if !ok {
					c.log.Error("channel closed")
				}
				if err := c.sendMessage(&data); err != nil {
					c.log.Error("send message")
				}
			case <- c.done:
				return
		}
	}
	
}

func (c *SMTPClient) Stop() {
	const op = "transport.smtp.mock.Stop"
	log := c.log.With(slog.String("op", op))
	log.Info("stop mock SMTP client")

	defer close(c.done)
}

func (c *SMTPClient) sendMessage(msg *model.Payload) error {
	const op = "transport.smtp.mock.sendMessage"
	log := c.log.With(slog.String("op", op))
	log.With(slog.String("email", msg.Email)).Info("send message")

	email, payload, err := c.prepareMessage(msg)
	if err != nil {
		log.Error("error prepare message", slerr.Err(err))
		return err
	}

	c.auth.With(
		slog.Attr(slog.String("to", email)),
		slog.Attr(slog.String("from", c.from)),
		slog.Attr(slog.String("payload", payload)),
	).Info(
		"message sent",
	)

	return nil
}

func (c *SMTPClient) prepareMessage(msg *model.Payload) (email, payload string, err error) {
	const op = "transport.smtp.mock.prepareMessage"
	log := c.log.With(slog.String("op", op))
	log.With(slog.String("email", msg.Email)).Info("prepare message")
	
	if msg.Email == EmptyString {
		log.Error("error prepare email", slerr.Err(ErrEmptyEmail))
		return "", "", fmt.Errorf("%s: %w", op, ErrEmptyEmail)
	}

	if msg.Message == EmptyString {
		log.Error("error prepare message", slerr.Err(ErrEmptyMessage))
		return "", "", fmt.Errorf("%s: %w", op, ErrEmptyMessage)
	}

	message := fmt.Sprintf(
		"To: %s\r\n"+
		"Subject: %s\r\n\r\n"+
		"You can use this password to log in: %s\nDo not forget to edit your password\r\n",
		 msg.Email,
		 msg.Header,
		 msg.Message,
	)

	return msg.Email, message, nil
}


