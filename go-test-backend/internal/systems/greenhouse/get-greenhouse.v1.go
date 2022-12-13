package greenhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//GetUser
type GetGreenhouse struct {
	GUID    uint64 `db:"guid"`
	Address string `db:"address"`
	Zip     uint64 `db:"zip"`
}

//GetUserV1Params
type GetGreenhouseV1Params struct {
	UUID uint64 `json:"uuid"`
	GUID uint64 `json:"guid"`
}

//GetUserV1Result
type GetGreenhouseV1Result struct {
	Greenhouse GetGreenhouse `json:"greenhouse"`
}

//GetUserV1 gets user by uuid
func (l *Greenhouse) GetGreenhouseV1(ctx context.Context, p *GetGreenhouseV1Params) (*GetGreenhouseV1Result, error) {
	greenhouse := GetGreenhouse{}
	v := ctx.Value("greenhouse_id")
	err := l.dbh.Get(&greenhouse, "SELECT guid, address, zip FROM greenhouses WHERE guid=?;", p.GUID)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	fmt.Println("context greenhouse_id", v, greenhouse)

	return &GetGreenhouseV1Result{Greenhouse: greenhouse}, nil
}

//GetUserHandler handles get user request
func (l *Greenhouse) GetGreenhouseHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	fmt.Println("get greenhouse vars", vars["guid"], vars["uuid"])
	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)

	req := &GetGreenhouseV1Params{
		UUID: uuid,
		GUID: guid,
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
