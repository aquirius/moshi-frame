package notification

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

// GetUserV1Params
type GetNotificationsV1Params struct {
	UUID uint64 `json:"uuid"`
}

// GetUserV1Result
type GetNotificationsV1Result struct {
	Notifications []GetNotification `json:"notifications"`
}

// GetUserV1 gets user by uuid
func (l *Notifications) GetNotificationsV1(ctx context.Context, p *GetNotificationsV1Params) (*GetNotificationsV1Result, error) {
	notifications := []GetNotification{}
	user := user.NewUserProvider(ctx, l.dbh, l.rdb, "")
	userID := user.User.GetUserID(p.UUID)
	err := l.dbh.Select(&notifications, "SELECT nuid, created_ts, checked_ts, done_ts, title, message FROM notifications WHERE user_id=?;", userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &GetNotificationsV1Result{Notifications: notifications}, nil
}

// GetUserHandler handles get user request
func (l *Notifications) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)

	var err error
	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		_, err := l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
	}

	req := &GetNotificationsV1Params{
		UUID: uuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetNotificationsV1(ctx, req)
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
