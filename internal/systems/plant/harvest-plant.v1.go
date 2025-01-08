package plant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// HarvestPlant
type HarvestPlant struct {
	PLUID string `db:"pluid"`
}

// HarvestPlantV1Params
type HarvestPlantV1Params struct {
	PUID uint64 `json:"puid"`
}

// HarvestPlantV1Result
type HarvestPlantV1Result struct {
}

// GetUserV1 gets user by uuid
func (l *Plant) HarvestPlantV1(ctx context.Context, p *HarvestPlantV1Params) (*HarvestPlantV1Result, error) {
	potID := l.GetPotID(p.PUID)

	query := "UPDATE plants SET harvested_ts=? WHERE pot_id=?;"
	_, err := l.dbh.Exec(query, time.Now().Unix(), potID)
	if err != nil {
		return nil, err
	}

	return &HarvestPlantV1Result{}, nil
}

// GetUserHandler handles get user request
func (l *Plant) HarvestPlantHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var err error
	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		//todoUPDATE plants SET harvested_ts=? WHERE pluid=?;
		_, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
	}
	req := &HarvestPlantV1Params{}

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.HarvestPlantV1(ctx, req)
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
