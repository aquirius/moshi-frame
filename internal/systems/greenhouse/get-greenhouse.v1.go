package greenhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUser
type GetGreenhouse struct {
	GUID    uint64 `db:"guid"`
	Address string `db:"address"`
	Zip     uint64 `db:"zip"`

	DisplayName *string  `db:"display_name"`
	Status      *string  `db:"status"`
	Destination *string  `db:"destination"`
	TempIn      *float64 `db:"tempIn"`
	TempOut     *float64 `db:"tempOut"`
	Humidity    *float64 `db:"humidity"`
	Brightness  *float64 `db:"brightness"`
	Co2         *float64 `db:"co2"`
}

// GetUserV1Params
type GetGreenhouseV1Params struct {
	UUID uint64 `json:"uuid"`
	GUID uint64 `json:"guid"`
}

// GetUserV1Result
type GetGreenhouseV1Result struct {
	Greenhouse GetGreenhouse `json:"greenhouse"`
}

// GetUserV1 gets user by uuid
func (l *Greenhouse) GetGreenhouseV1(ctx context.Context, p *GetGreenhouseV1Params) (*GetGreenhouseV1Result, error) {
	greenhouse := GetGreenhouse{}
	//v := ctx.Value("greenhouse_id")
	err := l.dbh.Get(&greenhouse, "SELECT guid, display_name, address, zip, status, destination, tempIn, tempOut, humidity, brightness, co2 FROM greenhouses WHERE guid=?;", p.GUID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &GetGreenhouseV1Result{Greenhouse: greenhouse}, nil
}

// GetUserHandler handles get user request
func (l *Greenhouse) GetGreenhouseHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)

	req := &GetGreenhouseV1Params{
		UUID: uuid,
		GUID: guid,
	}
	reqBody, _ := io.ReadAll(r.Body)
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
