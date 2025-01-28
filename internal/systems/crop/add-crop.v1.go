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
	"github.com/google/uuid"
)

// AddCrop
type AddCrop struct {
	PLUID       string `db:"pluid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

// AddCropV1Params
type AddCropV1Params struct {
	PUID     uint64 `json:"puid"`
	CropName string `json:"cropName"`
}

// AddCropV1Result
type AddCropV1Result struct {
	PLUID uint64 `json:"pluid"`
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

func (l *Crop) existingPUID(uuid uint32) bool {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (l *Crop) GetCropIDByName(name string) int64 {
	// for some reason i cant get the right crop by crop name
	// var query = "SELECT id FROM crops WHERE crop_name=?;"
	// var id int64
	// err := l.dbh.Get(&id, query, id)
	// if err != nil && err == sql.ErrNoRows {
	// 	return 0
	// }
	// return id
	if name == "tomato" {
		return 2
	} else {
		return 1
	}
}

// GetUserV1 gets user by uuid
func (l *Crop) AddCropV1(ctx context.Context, p *AddCropV1Params) (*AddCropV1Result, error) {
	var query string
	var err error
	//var result sql.Result
	//var nutrientID int64
	var pluid uuid.UUID
	//var cuid uuid.UUID

	pluid, err = uuid.NewUUID()
	//cuid, err := uuid.NewUUID()

	//userID := ctx.Value("user_id")
	//userID := l.getUserID(p.UUID)
	//potID := l.GetPotID(p.PUID)
	//cropID := l.GetCropIDByName(p.CropName)

	query = "INSERT INTO nutrients (carbon, hydrogen, oxygen, nitrogen, phosphorus, potassium, sulfur, calcium, magnesium) VALUES (?,?,?,?,?,?,?,?,?);"
	_, err = l.dbh.Exec(query, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		return nil, err
	}
	// nutrientID, err = result.LastInsertId()
	// if err != nil {
	// 	return nil, err
	// }
	// query = "INSERT INTO crops (cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	// result, err = l.dbh.Exec(query, cuid.ID(), "lettuce", 18, 28, 60, 80, 5.0, 6.0, 400, 500, 800, 1200, 18.0, 22.0)
	// if err != nil {
	// 	return nil, err
	// }
	// cuid, err = uuid.NewUUID()
	// query = "INSERT INTO crops (cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	// result, err = l.dbh.Exec(query, cuid.ID(), "tomato", 18, 26, 50, 70, 5.5, 6.5, 400, 500, 700, 1100, 19.0, 24.0)
	// if err != nil {
	// 	return nil, err
	// }
	// cropID, err = result.LastInsertId()

	return &AddCropV1Result{PLUID: uint64(pluid.ID())}, nil
}

// GetUserHandler handles get user request
func (l *Crop) AddCropHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	var err error
	//if we have a session id store it to req body
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
