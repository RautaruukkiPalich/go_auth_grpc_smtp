package app

import (
	"log/slog"

	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/app/kafka"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/app/smtp"
	"github.com/rautaruukkipalich/go_auth_grpc_smtp/internal/config"
)

type App struct {
	SMTP *smtp.SMTPClient
	Kafka kafka.Consumer
}

func New(log *slog.Logger, cfg *config.Config) *App{
	kafka := kafka.New(log, cfg)
	smtp := smtp.New(log, cfg)

	return &App{
		SMTP: smtp,
		Kafka: kafka,
	}
}

func (a *App) Run() {
	msgch := a.Kafka.Run()
	a.SMTP.Run(msgch)
}

func (a *App) Stop() {
	a.Kafka.Stop()
	a.SMTP.Stop()
}


