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
)

// GetUser
type GetUser struct {
	UUID        string `db:"uuid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

// GetUsersV1Params
type GetUsersV1Params struct {
}

// GetUsersV1Result
type GetUsersV1Result struct {
	Users []GetUser `json:"users"`
}

// GetUsersV1 gets all users
func (l *Users) GetUsersV1(p *GetUsersV1Params) (*GetUsersV1Result, error) {
	users := []GetUser{}

	err := l.dbh.Select(&users, "SELECT uuid, registered_ts, display_name, first_name, last_name, email, birthday FROM users;")
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	return &GetUsersV1Result{Users: users}, nil
}

// GetUsersHandler get all users request
func (l *Users) GetUsersHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &GetUsersV1Params{}

	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()
	session := SessionUser{}

	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		var redisSession string
		var err error
		redisSession, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(redisSession), &session)
		if err != nil {
			return nil, err
		}
	}

	if !session.Authenticated {
		return nil, errors.New("not logged in")
	}

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetUsersV1(nil)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
