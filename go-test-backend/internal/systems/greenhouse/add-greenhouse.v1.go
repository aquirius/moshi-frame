package greenhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/user"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//GetUser
type AddGreenhouse struct {
	GUID uint64 `json:"guid"`
}

//GetUserV1Params
type AddGreenhouseV1Params struct {
	UUID   uint64 `json:"uuid"`
	Adress string `json:"adress"`
	Zip    int16  `json:"zip"`
}

//GetUserV1Result
type AddGreenhouseV1Result struct {
	Greenhouse AddGreenhouse `json:"greenhouse"`
}

//GetUserV1 gets user by uuid
func (l *Greenhouses) AddGreenhouseV1(ctx context.Context, p *AddGreenhouseV1Params) (*AddGreenhouseV1Result, error) {

	guid, err := uuid.NewUUID()
	uguid, err := uuid.NewUUID()
	user := user.NewUserProvider(ctx, l.dbh, l.rdb, "")
	userID := user.User.GetUserID(p.UUID)

	query := "INSERT INTO greenhouses (guid, address, zip) VALUES (?,?,?);"
	result, err := l.dbh.Exec(query, guid.ID(), p.Adress, p.Zip)
	if err != nil {
		return nil, err
	}
	greenhouseID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO users_greenhouses (uguid, user_id, greenhouse_id) VALUES (?,?,?);"
	_, err = l.dbh.Exec(query, uguid.ID(), userID, greenhouseID)
	if err != nil {
		return nil, err
	}
	res := AddGreenhouse{
		GUID: uint64(guid.ID()),
	}

	return &AddGreenhouseV1Result{Greenhouse: res}, nil
}

//GetUserHandler handles get user request
func (l *Greenhouses) AddGreenhouseHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

	fmt.Println("redisSession add greenhouse", vars["uuid"], redisSession)
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	zip := 88965
	req := &AddGreenhouseV1Params{
		UUID:   uuid,
		Adress: "Hinterdupfing",
		Zip:    int16(zip),
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	fmt.Println("add greenhouse req", req)

	res, err := l.AddGreenhouseV1(ctx, req)
	fmt.Println("add greenhouse err", res, err)

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
