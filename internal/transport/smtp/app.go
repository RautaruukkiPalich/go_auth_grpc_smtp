package smtp

import (
	"fmt"
	"log/slog"
	"net/smtp"

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
	auth *smtp.Auth	
	addr string
	from string
}

func New(log *slog.Logger, cfg *config.SMTPConfig) *SMTPClient {
	
	auth := smtp.PlainAuth("", cfg.From, cfg.Pass, cfg.Host)
	done := make(chan struct{})


	return &SMTPClient{
		log: log,
		done: done,
		auth: &auth,
		addr: cfg.Addr,
		from: cfg.From,
	}
}

func (c *SMTPClient) Run(msgch chan model.Payload) {
	const op = "app.smtp.app.Run"
	log := c.log.With(slog.String("op", op))
	log.Info("start SMTP client")

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
	const op = "app.smtp.app.Run"
	log := c.log.With(slog.String("op", op))
	log.Info("stop SMTP client")

	defer close(c.done)
}

func (c *SMTPClient) sendMessage(msg *model.Payload) error {
	const op = "app.smtp.app.sendMessage"
	log := c.log.With(slog.String("op", op))
	log.With(slog.String("email", msg.Email)).Info("send message")

	email, payload, err := c.prepareMessage(msg)
	if err != nil {
		log.Error("error prepare message", op, slerr.Err(err))
		return err
	}
	
	if err := smtp.SendMail(c.addr, *c.auth, c.from, []string{email}, []byte(payload)); err != nil {
		log.Error("error sending message", op, slerr.Err(err))
		return err
	}

	return nil
}

func (c *SMTPClient) prepareMessage(msg *model.Payload) (email, payload string, err error) {
	const op = "app.smtp.app.prepareMessage"
	log := c.log.With(slog.String("op", op))
	log.With(slog.String("email", msg.Email)).Info("prepare message")
	
	if msg.Email == EmptyString {
		log.Error("error prepare email", op, slerr.Err(ErrEmptyEmail))
		return "", "", fmt.Errorf("%s: %w", op, ErrEmptyEmail)
	}

	if msg.Message == EmptyString {
		log.Error("error prepare message", op, slerr.Err(ErrEmptyMessage))
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


