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

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_NETWORK"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
}

func main() {
	out := []string{}

	db, err := sqlx.Open("mysql", buildMySQL())
	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	res, err := db.Exec("select * from users;")
	fmt.Println(res, err)

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

	fmt.Println(buildMySQL())

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
