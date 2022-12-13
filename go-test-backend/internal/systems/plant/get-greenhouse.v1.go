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

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

//GetUser
type GetPlant struct {
	UUID        string `db:"uuid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

type Greenhouse struct {
}

//GetUser
type GetGreenhouse struct {
	GUID    uint64 `db:"guid"`
	Address string `db:"address"`
	Zip     uint64 `db:"zip"`
}

//GetUserV1Params
type GetGreenhouseV1Params struct {
	UUID uint64 `json:"uuid"`
}

//GetUserV1Result
type GetGreenhouseV1Result struct {
	Greenhouse []GetGreenhouse `json:"greenhouse"`
}

//GetUserV1 gets user by uuid
func (l *Plant) GetGreenhouseV1(ctx context.Context, p *GetGreenhouseV1Params) (*GetGreenhouseV1Result, error) {

	userID := l.getUserID(p.UUID)
	greenhouses := []uint64{}
	v := ctx.Value("user_id")
	err := l.dbh.Select(&greenhouses, "SELECT greenhouse_id FROM users_greenhouses WHERE user_id=?;", userID)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	fmt.Println("context user_id", v, userID, greenhouses)

	getGreenhouses := []GetGreenhouse{}
	for _, v := range greenhouses {
		res := []GetGreenhouse{}
		err = l.dbh.Select(&res, "SELECT guid, address, zip FROM greenhouses WHERE id=?;", v)
		fmt.Println(&getGreenhouses, &res)
		if err == sql.ErrNoRows {
			fmt.Println("no rows")
			return nil, err
		}
		getGreenhouses = append(getGreenhouses, res...)
	}

	fmt.Println("context user_id", v, userID, getGreenhouses)

	return &GetGreenhouseV1Result{Greenhouse: getGreenhouses}, nil
}

//GetUserHandler handles get user request
func (l *Plant) GetGreenhouseHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	req := &GetGreenhouseV1Params{
		UUID: uuid,
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetGreenhouseV1(ctx, req)
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
