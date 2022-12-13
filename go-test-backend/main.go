package main

import (
	"context"
	"net/http"
	"test-backend/m/v2/internal/systems/plant"
	"test-backend/m/v2/internal/systems/user"

	"test-backend/m/v2/server"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//Runtime points to our systems
type Runtime struct {
	db    *sqlx.DB
	rdb   *redis.Client
	user  *user.User
	users *user.Users
	plant *plant.Plant
}

//BuildRuntime initializes our systems
func BuildRuntime() Runtime {
	//init empty context
	context := context.Background()
	//init server
	serverProvider := server.NewServerProvider()
	server := serverProvider.NewServer()
	//init users
	usersProvider := user.NewUsersProvider(context, &server.Sql, &server.Redis, "sql")
	users := usersProvider.NewUsers()
	//init user
	userProvider := user.NewUserProvider(context, &server.Sql, &server.Redis, "sql")
	user := userProvider.NewUser()
	//init plant
	plantProvider := plant.NewPlantProvider(context, &server.Sql, &server.Redis, "sql")
	plant := plantProvider.NewPlant()

	return Runtime{
		db:    &server.Sql,
		rdb:   &server.Redis,
		user:  user,
		users: users,
		plant: plant,
	}
}

func main() {
	//setup runtime with http server
	rt := BuildRuntime()
	mux := mux.NewRouter()

	//setup routes with their handlers
	userH := rt.user
	usersH := rt.users
	plantH := rt.plant

	mux.HandleFunc("/login", userH.ServeHTTP)
	mux.HandleFunc("/register", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse", plantH.ServeHTTP)

	mux.HandleFunc("/users", usersH.ServeHTTP)

	http.ListenAndServe(":1234", mux)
}
