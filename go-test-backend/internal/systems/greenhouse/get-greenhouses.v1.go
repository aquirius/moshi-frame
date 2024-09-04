package greenhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/user"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// GetUser
type GetGreenhouses struct {
	GUID    uint64 `db:"guid"`
	Address string `db:"address"`
	Zip     uint64 `db:"zip"`

	DisplayName *string  `db:"display_name"`
	Status      *string  `db:"status"`
	Destination *string  `db:"destination"`
	TempIn      *float64 `db:"tempIn"`
	TempOut     *float64 `db:"tempOut"`
	Humidity    *float64 `db:"humidity"`
	Brightness  *float64 `db:"brightness"`
	Co2         *float64 `db:"co2"`
}

// GetUserV1Params
type GetGreenhousesV1Params struct {
	UUID uint64 `json:"uuid"`
}

// GetUserV1Result
type GetGreenhousesV1Result struct {
	Greenhouses []GetGreenhouses `json:"greenhouses"`
}

// GetUserV1 gets user by uuid
func (l *Greenhouses) GetGreenhousesV1(ctx context.Context, p *GetGreenhousesV1Params) (*GetGreenhousesV1Result, error) {
	user := user.NewUserProvider(ctx, l.dbh, l.rdb, "")
	userID := user.User.GetUserID(p.UUID)
	greenhouses := []uint64{}
	v := ctx.Value("greenhouse_id")
	err := l.dbh.Select(&greenhouses, "SELECT greenhouse_id FROM users_greenhouses WHERE user_id=?;", userID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	fmt.Println("get greenhouses", v, len(greenhouses))

	getGreenhouses := []GetGreenhouses{}
	for _, v := range greenhouses {
		res := []GetGreenhouses{}
		err = l.dbh.Select(&res, "SELECT guid, display_name, address, zip, status, destination, tempIn, tempOut, humidity, brightness, co2 FROM greenhouses WHERE id=?;", v)
		if err == sql.ErrNoRows {
			return nil, err
		}
		getGreenhouses = append(getGreenhouses, res...)
	}

	return &GetGreenhousesV1Result{Greenhouses: getGreenhouses}, nil
}

// GetUserHandler handles get user request
func (l *Greenhouses) GetGreenhousesHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var redisSession string
	var err error
	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		redisSession, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
	}

	fmt.Println("redisSession", vars["uuid"], redisSession)
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	req := &GetGreenhousesV1Params{
		UUID: uuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetGreenhousesV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Fatal("error in json")
		return nil, err
	}
	return jsonBytes, nil
}
