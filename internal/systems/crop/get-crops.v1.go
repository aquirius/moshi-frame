package crop

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
type GetCropsV1Params struct {
	UUID uint64 `json:"uuid"`
}

// GetUserV1Result
type GetCropsV1Result struct {
	Crops []GetCrop `json:"crops"`
}

// GetUserV1 gets user by uuid
func (l *Crops) GetCropsV1(ctx context.Context, p *GetCropsV1Params) (*GetCropsV1Result, error) {
	user := user.NewUserProvider(ctx, l.dbh, l.rdb, "")
	userID := user.User.GetUserID(p.UUID)
	crops := []GetCrop{}
	err := l.dbh.Select(&crops, "SELECT crop_id, nutrient_id, pluid, created_ts, Croped_ts, harvested_ts FROM crops WHERE user_id=?;", userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &GetCropsV1Result{Crops: crops}, nil
}

// GetUserHandler handles get user request
func (l *Crops) GetCropsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	req := &GetCropsV1Params{
		UUID: uuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetCropsV1(ctx, req)
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
