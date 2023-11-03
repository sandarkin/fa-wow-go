package server

import (
	"io"
	"math/rand"
	"net"
	"strings"

	"github.com/rs/zerolog"
)

type Quotations struct {
	quotes []string
}

func NewQuotations(quotesBytes []byte) (*Quotations, error) {
	var quotes = strings.Split(string(quotesBytes), "\n")
	return &Quotations{quotes: quotes}, nil
}

func (b *Quotations) GetRandQuote() string {
	i := rand.Intn(len(b.quotes))
	return b.quotes[i]
}

func (b *Quotations) ServeRequest(conn net.Conn, requestLog zerolog.Logger) {
	requestLog.Info().Msg("write response")
	q := b.GetRandQuote()
	r := strings.NewReader(q)
	_, err := io.Copy(conn, r)
	if err != nil {
		requestLog.Warn().Err(err).Msg("failed to write response")
	}
}
