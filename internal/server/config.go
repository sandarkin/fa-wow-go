package server

type Config struct {
	ListenAddress string
	Difficulty    byte
	ProofSize     int
}

func NewConfig() *Config {
	c := new(Config)
	c.ListenAddress = "0.0.0.0:9000"
	c.Difficulty = byte(23)
	c.ProofSize = 64
	return c
}
