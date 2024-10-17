package server

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func buildMySQL() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_NETWORK"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
}

func (server *Server) connectSQL() *sqlx.DB {
	log.Print("connecting with mysql...")
	db, err := sqlx.Open("mysql", buildMySQL())
	if err != nil {
		fmt.Println(err)
	}

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	return db
}
