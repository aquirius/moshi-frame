package crop

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"test-backend/m/v2/internal/systems/user"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// AddCrop
type AddCrop struct {
	CUID         uint64  `db:"cuid"`
	CropName     string  `db:"crop_name"`
	AirTempMin   float64 `db:"air_temp_min"`
	AirTempMax   float64 `db:"air_temp_max"`
	HumidityMin  float64 `db:"humidity_min"`
	HumidityMax  float64 `db:"humidity_max"`
	PHLevelMin   float64 `db:"ph_level_min"`
	PHLevelMax   float64 `db:"ph_level_max"`
	OrpMin       float64 `db:"orp_min"`
	OrpMax       float64 `db:"orp_max"`
	TdsMin       uint16  `db:"tds_min"`
	TdsMax       uint16  `db:"tds_max"`
	WaterTempMin float64 `db:"water_temp_min"`
	WaterTempMax float64 `db:"water_temp_max"`
}

// AddCropV1Params
type AddCropV1Params struct {
	UUID         uint64  `json:"uuid"`
	CropName     string  `json:"cropName"`
	AirTempMin   float64 `json:"airTempMin"`
	AirTempMax   float64 `json:"airTempMax"`
	HumidityMin  float64 `json:"humidityMin"`
	HumidityMax  float64 `json:"humidityMax"`
	PHLevelMin   float64 `json:"phLevelMin"`
	PHLevelMax   float64 `json:"phLevelMax"`
	OrpMin       float64 `json:"orpMin"`
	OrpMax       float64 `json:"orpMax"`
	TdsMin       uint16  `json:"tdsMin"`
	TdsMax       uint16  `json:"tdsMax"`
	WaterTempMin float64 `json:"waterTempMin"`
	WaterTempMax float64 `json:"waterTempMax"`
}

// AddCropV1Result
type AddCropV1Result struct {
	CUID uint64 `json:"cuid"`
}

func (l *Crop) existingCUID(uuid uint32) bool {
	var query = "SELECT id FROM crops WHERE cuid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

// GetUserV1 gets user by uuid
func (l *Crops) AddCropV1(ctx context.Context, p *AddCropV1Params) (*AddCropV1Result, error) {
	var query string
	var err error
	var result sql.Result
	var ucuid uuid.UUID
	var cuid uuid.UUID

	user := user.NewUserProvider(ctx, l.dbh, l.rdb, "")
	userID := user.User.GetUserID(p.UUID)

	ucuid, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	cuid, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO crops (cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	result, err = l.dbh.Exec(query, cuid.ID(), p.CropName, p.AirTempMin, p.AirTempMax, p.HumidityMin, p.HumidityMax, p.PHLevelMin, p.PHLevelMax, p.OrpMin, p.OrpMax, p.TdsMin, p.TdsMax, p.WaterTempMin, p.WaterTempMax)
	if err != nil {
		return nil, err
	}

	cropID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO users_crops (ucuid, user_id, crop_id) VALUES (?,?,?);"
	_, err = l.dbh.Exec(query, ucuid, userID, cropID)
	if err != nil {
		return nil, err
	}

	return &AddCropV1Result{CUID: uint64(cuid.ID())}, nil
}

// AddCropHandler handles add crop request
func (l *Crops) AddCropHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var err error
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
		//todo
		_, err = l.rdb.Get(ctx, cookie.Value).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
	}
	req := &AddCropV1Params{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.AddCropV1(ctx, req)
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
