package server

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	httpAddr          = getenv("MILON_ADDR", "127.0.0.1:3000")
	mysqlDSN          = getenv("MILON_MYSQL", "milon:milon@tcp(127.0.0.1:3311)/milon")
	urlPrefixBackend  = getenv("MILON_URL_PREFIX_BACKEND", "http://127.0.0.1:3006")
	urlPrefixFrontend = getenv("MILON_URL_PREFIX_FRONTEND", "http://127.0.0.1:3000")
	fsPath            = getenv("MILON_FS_PATH", "uploads")
)

func getenv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func (server *Server) connectSQL() *sqlx.DB {
	log.Print("connecting with mysql...")
	db, err := sqlx.Open("mysql", mysqlDSN)
	if err != nil {
		fmt.Println(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
