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

func (l *Plants) GetPlantNutrients(nutriendID int) Nutrients {
	var query = "SELECT carbon, hydrogen, oxygen, nitrogen, phosphorus, potassium, sulfur, calcium, magnesium FROM nutrients WHERE id=?;"
	nutrients := Nutrients{}
	err := l.dbh.Get(&nutrients, query, nutriendID)
	if err != nil && err == sql.ErrNoRows {
		return Nutrients{0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	return nutrients
}

func (l *Plants) GetPlantCrop(cropID int) Crop {
	var query = "SELECT cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max FROM crops WHERE id=?;"
	crop := Crop{}
	err := l.dbh.Get(&crop, query, cropID)
	if err != nil && err == sql.ErrNoRows {
		return Crop{0, "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	return crop
}

// GetUserV1 gets user by uuid
func (l *Plants) GetPlantsV1(ctx context.Context, p *GetPlantsV1Params) (*GetPlantsV1Result, error) {
	plants := []GetPlant{}
	for _, v := range p.PUID {
		potID := l.GetPotID(v)
		res := GetPlant{}
		err := l.dbh.Get(&res, "SELECT crop_id, nutrient_id, pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;", potID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		res.Nutrients = l.GetPlantNutrients(res.NutrientID)
		res.Crop = l.GetPlantCrop(res.CropID)

		plants = append(plants, res)
	}

	return &GetPlantsV1Result{Plants: plants}, nil
}

// GetUserHandler handles get user request
func (l *Plants) GetPlantsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
