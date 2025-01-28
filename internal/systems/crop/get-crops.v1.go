package crop

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	redis "github.com/go-redis/redis/v8"
)

// GetUserV1Params
type GetCropsV1Params struct {
	PUID []uint64 `json:"puid"`
}

// GetUserV1Result
type GetCropsV1Result struct {
	Crops []GetCrop `json:"crops"`
}

// GetUserV1 gets user by uuid
func (l *Crops) GetCropsV1(ctx context.Context, p *GetCropsV1Params) (*GetCropsV1Result, error) {
	crops := []GetCrop{}
	err := l.dbh.Select(&crops, "SELECT crop_id, nutrient_id, pluid, created_ts, Croped_ts, harvested_ts FROM crops WHERE pot_id=?;", p.PUID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &GetCropsV1Result{Crops: crops}, nil
}

// GetUserHandler handles get user request
func (l *Crops) GetCropsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	//vars := mux.Vars(r)
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

	req := &GetCropsV1Params{}
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
