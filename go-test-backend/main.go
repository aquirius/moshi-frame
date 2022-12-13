package main

import (
	"context"
	"fmt"
	"net/http"
	"test-backend/m/v2/internal/systems/greenhouse"
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
	db          *sqlx.DB
	rdb         *redis.Client
	user        *user.User
	users       *user.Users
	plant       *plant.Plant
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
	greenhouseH := rt.greenhouse
	greenhousesH := rt.greenhouses

	mux.HandleFunc("/login", userH.ServeHTTP)
	mux.HandleFunc("/register", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse", greenhousesH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}", greenhouseH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/get-stacks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		switch {
		case r.Method == http.MethodGet:
			fmt.Println("get stacks")
			res, err := rt.plant.GetStacksHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		default:
			return
		}
	})

	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/get-pots", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		switch {
		case r.Method == http.MethodPost:
			method := r.Header.Get("Method")
			var res []byte
			var err error
			if method == "add" {
				fmt.Println("post get stack add")
				res, err = rt.plant.GetPotsHandler(w, r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(err.Error()))
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		default:
			return
		}
	})
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/add-stack", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		switch {
		case r.Method == http.MethodPost:
			method := r.Header.Get("Method")
			var res []byte
			var err error
			if method == "add" {
				fmt.Println("post stack add")
				res, err = rt.plant.AddStackHandler(w, r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(err.Error()))
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		default:
			return
		}
	})

	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/add-pot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		switch {
		case r.Method == http.MethodPost:
			method := r.Header.Get("Method")
			var res []byte
			var err error
			if method == "add" {
				fmt.Println("post stack add")
				res, err = rt.plant.AddPotHandler(w, r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(err.Error()))
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		default:
			return
		}
	})

	mux.HandleFunc("/users", usersH.ServeHTTP)

	http.ListenAndServe(":1234", mux)
}
