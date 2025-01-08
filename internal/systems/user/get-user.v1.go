package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// GetUserV1Params
type GetUserV1Params struct {
	UUID string `json:"uuid"`
}

// GetUserV1Result
type GetUserV1Result struct {
	User GetUser `json:"user"`
}

// GetUserV1 gets user by uuid
func (l *User) GetUserV1(ctx context.Context, p *GetUserV1Params) (*GetUserV1Result, error) {
	user := GetUser{}
	ctx.Value("session-id")

	err := l.dbh.Get(&user, "SELECT uuid, registered_ts, display_name, first_name, last_name, email, birthday FROM users WHERE uuid=?", p.UUID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &GetUserV1Result{User: user}, nil
}

// GetUserHandler handles get user request
func (l *User) GetUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
	} else {
		return nil, errors.New("not logged in")
	}

	fmt.Println("redisSession", redisSession)

	req := &GetUserV1Params{
		UUID: vars["uuid"],
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetUserV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
