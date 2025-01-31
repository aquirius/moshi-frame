package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/greenhouse"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetUser
type AddStack struct {
	SUID uint64 `db:"suid"`
}

// GetUserV1Params
type AddStackV1Params struct {
	GUID uint64 `json:"guid"`
}

// GetUserV1Result
type AddStackV1Result struct {
	Stack AddStack `json:"stack"`
}

// GetUserV1 gets user by uuid
func (l *Stack) AddStackV1(ctx context.Context, p *AddStackV1Params) (*AddStackV1Result, error) {

	suid, err := uuid.NewUUID()
	greenhouse := greenhouse.NewGreenhouseProvider(ctx, l.dbh, l.rdb, "")
	greenhouseID := greenhouse.Greenhouse.GetGreenhouseID(p.GUID)

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

// GetUserHandler handles get user request
func (l *Stack) AddStackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var err error
	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		_, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
	}

	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)
	req := &AddStackV1Params{
		GUID: guid,
	}
	reqBody, _ := io.ReadAll(r.Body)
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
