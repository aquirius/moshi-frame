package main

import (
	"context"
	"net/http"
	greenhouse "test-backend/m/v2/internal/systems/greenhouse"
	plant "test-backend/m/v2/internal/systems/plant"
	pot "test-backend/m/v2/internal/systems/pot"
	stack "test-backend/m/v2/internal/systems/stack"
	user "test-backend/m/v2/internal/systems/user"

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
	pot   *pot.Pot
	stack *stack.Stack

	greenhouse  *greenhouse.Greenhouse
	greenhouses *greenhouse.Greenhouses
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
	//init plant
	potProvider := pot.NewPotProvider(context, &server.Sql, &server.Redis, "sql")
	pot := potProvider.NewPot()
	//init plant
	stackProvider := stack.NewStackProvider(context, &server.Sql, &server.Redis, "sql")
	stack := stackProvider.NewStack()
	//init greenhouses
	greenhousesProvider := greenhouse.NewGreenhousesProvider(context, &server.Sql, &server.Redis, "sql")
	greenhouses := greenhousesProvider.NewGreenhouses()
	//init greenhouse
	greenhouseProvider := greenhouse.NewGreenhouseProvider(context, &server.Sql, &server.Redis, "sql")
	greenhouse := greenhouseProvider.NewGreenhouse()

	return Runtime{
		db:          &server.Sql,
		rdb:         &server.Redis,
		user:        user,
		users:       users,
		plant:       plant,
		pot:         pot,
		stack:       stack,
		greenhouse:  greenhouse,
		greenhouses: greenhouses,
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
	potH := rt.pot
	stackH := rt.stack
	greenhouseH := rt.greenhouse
	greenhousesH := rt.greenhouses

	mux.HandleFunc("/login", userH.ServeHTTP)
	mux.HandleFunc("/logout", userH.ServeHTTP)
	mux.HandleFunc("/register", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouses", greenhousesH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}", greenhouseH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/plant", plantH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/pot", potH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack", stackH.ServeHTTP)
	mux.HandleFunc("/users", usersH.ServeHTTP)
	//todo .env
	http.ListenAndServe("127.0.1:1234", mux)
}
