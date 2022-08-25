package server

import (
	redis "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

//ServerProvider provides *Server
type ServerProvider struct {
	Server *Server
}

//Server
type Server struct {
	Sql   sqlx.DB
	Redis redis.Client
}

//NewServerProvider returns a new Server provider
func NewServerProvider() *ServerProvider {
	return &ServerProvider{
		&Server{},
	}
}

//NewServer initializes all databases
func (b *ServerProvider) NewServer() *Server {
	b.Server.Sql = *b.Server.connectSQL()
	b.Server.Redis = *redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "milon",
		DB:       1,
	})
	return b.Server
}
