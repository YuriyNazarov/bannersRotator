package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	internalamqp "github.com/YuriyNazarov/bannersRotator/internal/amqp/rabbit"
	internalapp "github.com/YuriyNazarov/bannersRotator/internal/app"
	internalconfig "github.com/YuriyNazarov/bannersRotator/internal/config"
	"github.com/YuriyNazarov/bannersRotator/internal/logger"
	internalselector "github.com/YuriyNazarov/bannersRotator/internal/selector"
	internalserver "github.com/YuriyNazarov/bannersRotator/internal/server"
	internalstorage "github.com/YuriyNazarov/bannersRotator/internal/storage"
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

	storage := internalstorage.New(logg, config.Database.DSN)
	defer storage.Close()

	selector := internalselector.New()
	app := internalapp.New(logg, storage, storage, selector, amqp)

	server := internalserver.NewServer(logg, app, net.JoinHostPort(config.Server.Host, config.Server.Port))
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop server: " + err.Error())
		}
	}()

	logg.Info("UP")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start server: " + err.Error())
		cancel()
	}
}
