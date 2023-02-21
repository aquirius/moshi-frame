package plant

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//GetUser
type AddPlant struct {
	UUID        string `db:"uuid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

//GetUserV1Params
type AddPlantV1Params struct {
	UUID uint64 `json:"uuid"`
}

//GetUserV1Result
type AddPlantV1Result struct {
	UUID uint64 `json:"uuid"`
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

//GetUserV1 gets user by uuid
func (l *Plant) AddPlantV1(ctx context.Context, p *AddPlantV1Params) (*AddPlantV1Result, error) {
	pluid, err := uuid.NewUUID()
	puid, err := uuid.NewUUID()

	//userID := ctx.Value("user_id")

	//userID := l.getUserID(p.UUID)

	query := "INSERT INTO nutrients (carbon, hydrogen, oxygen, nitrogen, phosphorus, potassium, sulfur, calcium, magnesium) VALUES (?,?,?,?,?,?,?,?,?);"
	result, err := l.dbh.Exec(query, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		return nil, err
	}
	nutrientID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO pots (puid, stack_id, user_id) VALUES (?,?,?);"
	result, err = l.dbh.Exec(query, puid.ID(), 0, 1)
	if err != nil {
		return nil, err
	}
	potID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO plants (pluid, created_ts, harvested_ts, nutrient_id, pot_id) VALUES (?,?,?,?,?);"
	_, err = l.dbh.Exec(query, pluid.ID(), time.Now().Unix(), time.Now().Unix(), nutrientID, potID)
	if err != nil {
		return nil, err
	}

	return &AddPlantV1Result{UUID: uint64(pluid.ID())}, nil
}

//GetUserHandler handles get user request
func (l *Plant) AddPlantHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
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
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)

	req := &AddPlantV1Params{
		UUID: uuid,
	}
	fmt.Println(req)

	reqBody, _ := ioutil.ReadAll(r.Body)
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
