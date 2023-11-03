package server

import (
	"errors"
	"net"

	"github.com/rs/zerolog"

	"github.com/sandarkin/fa-wow-go/internal/pow"
)

type Server struct {
	conf       *Config
	log        zerolog.Logger
	sock       net.Listener
	handler    func(net.Conn, zerolog.Logger)
	powReceive pow.Receiver
}

func StartServer(conf *Config, log zerolog.Logger, handler func(net.Conn, zerolog.Logger)) (*Server, error) {
	socket, err := net.Listen("tcp", conf.ListenAddress)
	if err != nil {
		return nil, err
	}
	s := &Server{
		conf:       conf,
		log:        log,
		sock:       socket,
		handler:    handler,
		powReceive: pow.NewReceiver(conf.Difficulty, conf.ProofSize),
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	return s.sock.Close()
}

func (s *Server) listen() {
	for i := 0; ; i++ {
		conn, err := s.sock.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			s.log.Warn().Err(err).Msg("failed to listen socket")
			continue
		}
		go s.serveConn(conn, i)
	}
}

func (s *Server) serveConn(conn net.Conn, connID int) {
	defer conn.Close()

	log := s.log.With().
		Int("id", connID).
		Str("addr", conn.RemoteAddr().String()).
		Logger()
	log.Trace().Msg("receive connection")

	checkDuration, err := s.powReceive(conn)
	if err != nil {
		log.Warn().Err(err).Dur("check_duration", checkDuration).Msg("refuse connection")
		return
	}
	log.Debug().
		Int("difficulty", int(s.conf.Difficulty)).
		Dur("check_duration", checkDuration).
		Msg("is valid proof")

	s.handler(conn, log)
}
