package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// GetUserV1Params
type GetUserSettingsV1Params struct {
	UUID string `json:"uuid"`
}

// GetUserV1Result
type GetUserSettingsV1Result struct {
	User GetUser `json:"user"`
}

type RedisSession struct {
	UUID          string `json:"UUID"`
	Password      string `json:"Password"`
	Authenticated bool   `json:"Authenticated"`
}

// GetUserV1 gets user by uuid
func (l *User) GetUserSettingsV1(ctx context.Context, p *GetUserSettingsV1Params) (*GetUserSettingsV1Result, error) {
	user := GetUser{}
	//v := ctx.Value("session-id")
	err := l.dbh.Get(&user, "SELECT uuid, registered_ts, display_name, first_name, last_name, email, birthday FROM users WHERE uuid=?", p.UUID)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	return &GetUserSettingsV1Result{User: user}, nil
}

// GetUserHandler handles get user request
func (l *User) GetUserSettingsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var redisSession string
	var redisJson RedisSession
	var err error
	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		redisSession, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			return nil, errors.New("not logged in")
		} else if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("not logged in")
	}

	err = json.Unmarshal([]byte(redisSession), &redisJson)
	if err != nil {
		redisJson.Authenticated = false
		return nil, errors.New("you motherfucker")
	}

	if !redisJson.Authenticated {
		return nil, errors.New("not authenticated")
	}

	req := &GetUserSettingsV1Params{
		UUID: vars["uuid"],
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetUserSettingsV1(ctx, req)
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
