package main

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
	return fmt.Sprintf("root:%s@%s(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_NETWORK"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func main() {
	out := []string{}

	db, err := sqlx.Open("mysql", buildMySQL())
	if err != nil {
		panic(err.Error())
	}

	greenhouses, err := os.ReadFile("./schemas/greenhouses.sql")
	if err != nil {
		panic(err.Error())
	}
	userGreenhouse, err := os.ReadFile("./schemas/users-greenhouses.sql")
	if err != nil {
		panic(err.Error())
	}
	notifications, err := os.ReadFile("./schemas/notifications.sql")
	if err != nil {
		panic(err.Error())
	}
	crops, err := os.ReadFile("./schemas/crops.sql")
	if err != nil {
		panic(err.Error())
	}
	nutrients, err := os.ReadFile("./schemas/nutrients.sql")
	if err != nil {
		panic(err.Error())
	}
	plans, err := os.ReadFile("./schemas/plans.sql")
	if err != nil {
		panic(err.Error())
	}
	plants, err := os.ReadFile("./schemas/plants.sql")
	if err != nil {
		panic(err.Error())
	}
	pots, err := os.ReadFile("./schemas/pots.sql")
	if err != nil {
		panic(err.Error())
	}
	stacks, err := os.ReadFile("./schemas/stacks.sql")
	if err != nil {
		panic(err.Error())
	}
	users, err := os.ReadFile("./schemas/users.sql")
	if err != nil {
		panic(err.Error())
	}
	sprouts, err := os.ReadFile("./schemas/sprouts.sql")
	if err != nil {
		panic(err.Error())
	}

	out = append(out,
		string(greenhouses),
		string(notifications),
		string(crops),
		string(nutrients),
		string(plans),
		string(plants),
		string(pots),
		string(stacks),
		string(users),
		string(userGreenhouse),
		string(sprouts))

	for _, v := range out {
		res, err := db.Exec(v)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}
}
