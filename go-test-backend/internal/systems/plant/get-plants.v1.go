package plant

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// GetUserV1Params
type GetPlantsV1Params struct {
	PUID []uint64 `json:"puid"`
}

// GetUserV1Result
type GetPlantsV1Result struct {
	Plants []GetPlant `json:"plants"`
}

// GetUserV1 gets user by uuid
func (l *Plant) GetPlantsV1(ctx context.Context, p *GetPlantsV1Params) (*GetPlantsV1Result, error) {
	plants := []GetPlant{}
	for _, v := range p.PUID {
		res := GetPlant{}
		err := l.dbh.Select(&res, "SELECT nutrient_id, pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;", v)
		if err == sql.ErrNoRows {
			fmt.Println("no rows")
			return nil, err
		}
		res.Nutrients = *l.GetPlantNutrients(res.NutrientID)
		plants = append(plants, res)
	}

	return &GetPlantsV1Result{Plants: plants}, nil
}

// GetUserHandler handles get user request
func (l *Plant) GetPlantsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

	fmt.Println("redisSession", vars["uuid"], redisSession)

	req := &GetPlantsV1Params{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetPlantsV1(ctx, req)
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
