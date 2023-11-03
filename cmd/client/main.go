package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/sandarkin/fa-wow-go/internal/client"
)

func main() {
	zerolog.DurationFieldUnit = time.Millisecond

	ctx, ctxCancel := context.WithCancel(context.TODO())
	defer ctxCancel()

	conf := client.NewConfig()

	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().Timestamp().
		Logger()

	log.Debug().
		Str("server_addr", conf.ServerAddr).
		Int("workers", conf.Workers).
		Dur("timeout", conf.Timeout).
		Msg("server client")

	client.StartWorkers(ctx, conf, log)

	waitForExit()
}

func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
