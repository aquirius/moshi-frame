package server

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func buildMySQL() string {
	return fmt.Sprintf("root:%s@%s(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_NETWORK"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func (server *Server) connectSQL() *sqlx.DB {
	log.Print("connecting with mysql...")
	db, err := sqlx.Open("mysql", buildMySQL())
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	return db
}
