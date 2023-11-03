package main

import (
	_ "embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/sandarkin/fa-wow-go/internal/server"
)

//go:embed quotes.txt
var quotesBytes []byte

func main() {
	zerolog.DurationFieldUnit = time.Millisecond

	conf := server.NewConfig()

	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().Timestamp().
		Logger()

	log.Debug().
		Str("listen_address", conf.ListenAddress).
		Int("proof_size", conf.ProofSize).
		Int("difficulty", int(conf.Difficulty)).
		Msg("server started")

	quotations, err := server.NewQuotations(quotesBytes)
	check(log, err)

	server, err := server.StartServer(conf, log, quotations.ServeRequest)
	check(log, err)
	defer server.Close()

	waitForExit()
}

func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func check(log zerolog.Logger, err error) {
	if err != nil {
		log.Fatal().Err(err).CallerSkipFrame(1).Msg("start failed")
		panic(err)
	}
}
