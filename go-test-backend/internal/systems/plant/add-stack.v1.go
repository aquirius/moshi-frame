package plant

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//GetUser
type AddStack struct {
	SUID uint64 `json:"guid"`
}

//GetUserV1Params
type AddStackV1Params struct {
	GUID uint64 `json:"guid"`
}

//GetUserV1Result
type AddStackV1Result struct {
	Stack AddStack `json:"stack"`
}

//GetUserV1 gets user by uuid
func (l *Plant) AddStackV1(ctx context.Context, p *AddStackV1Params) (*AddStackV1Result, error) {

	suid, err := uuid.NewUUID()
	greenhouseID := l.getGreenhouseID(p.GUID)

	query := "INSERT INTO stacks (suid, greenhouse_id) VALUES (?,?);"
	_, err = l.dbh.Exec(query, suid.ID(), greenhouseID)
	if err != nil {
		return nil, err
	}

	res := &AddStack{
		SUID: uint64(suid.ID()),
	}

	return &AddStackV1Result{Stack: *res}, nil
}

//GetUserHandler handles get user request
func (l *Plant) AddStackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

	fmt.Println("redisSession", vars["guid"], redisSession)
	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)
	req := &AddStackV1Params{
		GUID: guid,
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.AddStackV1(ctx, req)
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
