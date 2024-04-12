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
	log *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App{
	kafka := kafka.New(log, &cfg.Kafka)
	smtp := smtp.New(log, &cfg.SMTP)

	return &App{
		SMTP: smtp,
		Kafka: kafka,
		log: log,
	}
}

func (a *App) Run() {
	const op = "app.app.Run"
	log := a.log.With(slog.String("op", op))
	log.Info("starting app")

	msgch := a.Kafka.Run()
	a.SMTP.Run(msgch)
}

func (a *App) Stop() {
	const op = "app.app.Run"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping app")

	defer a.Kafka.Stop()
	defer a.SMTP.Stop()
}


