package plant

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// AddPlant
type AddPlant struct {
	PLUID       string `db:"pluid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

// AddPlantV1Params
type AddPlantV1Params struct {
	PUID uint64 `json:"puid"`
}

// AddPlantV1Result
type AddPlantV1Result struct {
	PLUID uint64 `json:"pluid"`
}

func (l *Plant) existingPLUID(uuid uint32) bool {
	var query = "SELECT id FROM plants WHERE pluid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (l *Plant) existingPUID(uuid uint32) bool {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

// GetUserV1 gets user by uuid
func (l *Plant) AddPlantV1(ctx context.Context, p *AddPlantV1Params) (*AddPlantV1Result, error) {
	pluid, err := uuid.NewUUID()

	//userID := ctx.Value("user_id")
	//userID := l.getUserID(p.UUID)
	potID := l.GetPotID(p.PUID)

	query := "INSERT INTO nutrients (carbon, hydrogen, oxygen, nitrogen, phosphorus, potassium, sulfur, calcium, magnesium) VALUES (?,?,?,?,?,?,?,?,?);"
	result, err := l.dbh.Exec(query, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		return nil, err
	}
	nutrientID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO plants (pluid, created_ts, planted_ts, nutrient_id, pot_id) VALUES (?,?,?,?,?);"
	result, err = l.dbh.Exec(query, pluid.ID(), time.Now().Unix(), time.Now().Unix(), nutrientID, potID)
	if err != nil {
		return nil, err
	}
	plantID, err := result.LastInsertId()
	fmt.Println("last inser id of plant", plantID)

	return &AddPlantV1Result{PLUID: uint64(pluid.ID())}, nil
}

// GetUserHandler handles get user request
func (l *Plant) AddPlantHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
	req := &AddPlantV1Params{}
	fmt.Println(req)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.AddPlantV1(ctx, req)
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
