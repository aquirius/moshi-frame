package plant

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/pot"

	"github.com/gorilla/mux"
)

// GetPlant
type GetPlant struct {
	PLUID       uint64 `db:"pluid"`
	CreatedTS   uint64 `db:"created_ts"`
	PlantedTS   uint64 `db:"planted_ts"`
	HarvestedTS uint64 `db:"harvested_ts"`
}

// GetPlantV1Params
type GetPlantV1Params struct {
	PUID uint64 `json:"puid"`
}

// GetPlantV1Result
type GetPlantV1Result struct {
	Plant GetPlant `json:"plant"`
}

// GetPlantV1 gets pots by suid
func (l *Plant) GetPlantV1(ctx context.Context, p *GetPlantV1Params) (*GetPlantV1Result, error) {
	plant := GetPlant{}

	pot := pot.NewPotProvider(ctx, l.dbh, l.rdb, "")
	potID := pot.Pot.GetPotID(p.PUID)

	err := l.dbh.Get(&plant, "SELECT pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;", potID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &GetPlantV1Result{Plant: plant}, nil
}

// GetUserHandler handles get user request
func (l *Plant) GetPlantHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	puid, _ := strconv.ParseUint(vars["puid"], 0, 32)
	req := &GetPlantV1Params{
		PUID: puid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetPlantV1(ctx, req)
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
