package server

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	httpAddr          = getenv("SPROUT_ADDR", "127.0.0.1:3000")
	mysqlDSN          = getenv("SPROUT_MYSQL", "root:moshi@tcp(127.0.0.1:3311)/sprout")
	urlPrefixBackend  = getenv("SPROUT_URL_PREFIX_BACKEND", "http://127.0.0.1:3006")
	urlPrefixFrontend = getenv("SPROUT_URL_PREFIX_FRONTEND", "http://127.0.0.1:3000")
	fsPath            = getenv("SPROUT_FS_PATH", "uploads")
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

	return db
}
