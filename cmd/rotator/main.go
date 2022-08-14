package main

import (
	"context"
	internalamqp "github.com/YuriyNazarov/bannersRotator/internal/amqp"
	internalconfig "github.com/YuriyNazarov/bannersRotator/internal/config"
	"github.com/YuriyNazarov/bannersRotator/internal/logger"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	config, err := internalconfig.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	logg := logger.NewLogger(config.Logger)
	defer logg.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	amqp := internalamqp.NewRabbit(ctx, logg, config.Queue)
}
