package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	crop "test-backend/m/v2/internal/systems/crop"
	greenhouse "test-backend/m/v2/internal/systems/greenhouse"
	notification "test-backend/m/v2/internal/systems/notification"
	plant "test-backend/m/v2/internal/systems/plant"
	pot "test-backend/m/v2/internal/systems/pot"
	sprout "test-backend/m/v2/internal/systems/sprout"
	stack "test-backend/m/v2/internal/systems/stack"
	user "test-backend/m/v2/internal/systems/user"

	server "test-backend/m/v2/server"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Runtime points to our systems
type Runtime struct {
	db  *sqlx.DB
	rdb *redis.Client

	user  *user.User
	users *user.Users

	plant  *plant.Plant
	plants *plant.Plants

	pot   *pot.Pot
	stack *stack.Stack

	crop  *crop.Crop
	crops *crop.Crops

	greenhouse  *greenhouse.Greenhouse
	greenhouses *greenhouse.Greenhouses

	sprout *sprout.Sprout

	notification  *notification.Notification
	notifications *notification.Notifications
}

// BuildRuntime initializes our systems
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
	//init plants
	plantsProvider := plant.NewPlantsProvider(context, &server.Sql, &server.Redis, "sql")
	plants := plantsProvider.NewPlants()
	//init plant
	plantProvider := plant.NewPlantProvider(context, &server.Sql, &server.Redis, "sql")
	plant := plantProvider.NewPlant()
	//init crops
	cropsProvider := crop.NewCropsProvider(context, &server.Sql, &server.Redis, "sql")
	crops := cropsProvider.NewCrops()
	//init crop
	cropProvider := crop.NewCropProvider(context, &server.Sql, &server.Redis, "sql")
	crop := cropProvider.NewCrop()
	//init pot
	potProvider := pot.NewPotProvider(context, &server.Sql, &server.Redis, "sql")
	pot := potProvider.NewPot()
	//init stack
	stackProvider := stack.NewStackProvider(context, &server.Sql, &server.Redis, "sql")
	stack := stackProvider.NewStack()
	//init greenhouses
	greenhousesProvider := greenhouse.NewGreenhousesProvider(context, &server.Sql, &server.Redis, "sql")
	greenhouses := greenhousesProvider.NewGreenhouses()
	//init greenhouse
	greenhouseProvider := greenhouse.NewGreenhouseProvider(context, &server.Sql, &server.Redis, "sql")
	greenhouse := greenhouseProvider.NewGreenhouse()
	//init sprout
	sproutProvider := sprout.NewSproutProvider(context, &server.Sql, &server.Redis, "sql")
	sprout := sproutProvider.NewSprout()
	//init notifications
	notificationsProvider := notification.NewNotificationsProvider(context, &server.Sql, &server.Redis, "sql")
	notifications := notificationsProvider.NewNotifications()
	//init notification
	notificationProvider := notification.NewNotificationProvider(context, &server.Sql, &server.Redis, "sql")
	notification := notificationProvider.NewNotification()

	return Runtime{
		db:  &server.Sql,
		rdb: &server.Redis,

		user:  user,
		users: users,

		plant:  plant,
		plants: plants,

		crop:  crop,
		crops: crops,

		pot:           pot,
		stack:         stack,
		greenhouse:    greenhouse,
		greenhouses:   greenhouses,
		sprout:        sprout,
		notification:  notification,
		notifications: notifications,
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
	plantsH := rt.plants
	cropH := rt.crop
	cropsH := rt.crops
	potH := rt.pot
	stackH := rt.stack
	greenhouseH := rt.greenhouse
	greenhousesH := rt.greenhouses
	sproutH := rt.sprout
	notificationH := rt.notification
	notificationsH := rt.notifications
	mux.HandleFunc("/ws", sproutH.HandleWebSocket)
	mux.HandleFunc("/login", userH.ServeHTTP)
	mux.HandleFunc("/logout", userH.ServeHTTP)
	mux.HandleFunc("/register", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}", userH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/crops", cropsH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/crop/{cuid}", cropH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/notifications", notificationsH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/notification/{nuid}", notificationH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouses", greenhousesH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}", greenhouseH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack", stackH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack/{suid}", sproutH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack/{suid}/pot", potH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack/{suid}/plants", plantsH.ServeHTTP)
	mux.HandleFunc("/user/{uuid}/greenhouse/{guid}/stack/{suid}/pot/{puid}/plant", plantH.ServeHTTP)
	mux.HandleFunc("/users", usersH.ServeHTTP)
	//todo middleware
	http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("BACKEND_HOST"), os.Getenv("BACKEND_PORT")), mux)
}
