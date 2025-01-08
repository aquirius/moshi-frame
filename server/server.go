package server

import (
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// ServerProvider provides *Server
type ServerProvider struct {
	Server *Server
}

// Server
type Server struct {
	Sql   sqlx.DB
	Redis redis.Client
}

// NewServerProvider returns a new Server provider
func NewServerProvider() *ServerProvider {
	return &ServerProvider{
		&Server{},
	}
}

// NewServer initializes all databases
func (b *ServerProvider) NewServer() *Server {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	b.Server.Sql = *b.Server.connectSQL()
	b.Server.Redis = *redis.NewClient(&redis.Options{
		Username: os.Getenv("REDIS_USER"),
		Addr:     os.Getenv("REDIS_URL") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return b.Server
}
