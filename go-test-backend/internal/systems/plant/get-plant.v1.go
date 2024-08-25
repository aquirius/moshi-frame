package plant

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/pot"

	"github.com/gorilla/mux"
)

type Nutrients struct {
	Carbon     uint16 `db:"carbon"`
	Hydrogen   uint16 `db:"hydrogen"`
	Oxygen     uint16 `db:"oxygen"`
	Nitrogen   uint16 `db:"nitrogen"`
	Phosphorus uint16 `db:"phosphorus"`
	Potassium  uint16 `db:"potassium"`
	Sulfur     uint16 `db:"sulfur"`
	Calcium    uint16 `db:"calcium"`
	Magnesium  uint16 `db:"magnesium"`
}

// GetPlant
type GetPlant struct {
	NutrientID  int    `db:"nutrient_id"`
	PLUID       uint64 `db:"pluid"`
	CreatedTS   uint64 `db:"created_ts"`
	PlantedTS   uint64 `db:"planted_ts"`
	HarvestedTS uint64 `db:"harvested_ts"`

	Nutrients Nutrients
}

// GetPlantV1Params
type GetPlantV1Params struct {
	PUID uint64 `json:"puid"`
}

// GetPlantV1Result
type GetPlantV1Result struct {
	Plant GetPlant `json:"plant"`
}

func (l *Plant) GetPlantNutrients(nutriendID int) *Nutrients {
	var query = "SELECT carbon, hydrogen, oxygen, nitrogen, phosphorus, potassium, sulfur, calcium, magnesium FROM nutrients WHERE id=?;"
	nutrients := Nutrients{}
	err := l.dbh.Get(&nutrients, query, nutriendID)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	fmt.Println(&nutrients)
	return &nutrients
}

// GetPlantV1 gets pots by suid
func (l *Plant) GetPlantV1(ctx context.Context, p *GetPlantV1Params) (*GetPlantV1Result, error) {
	plant := GetPlant{}

	pot := pot.NewPotProvider(ctx, l.dbh, l.rdb, "")
	potID := pot.Pot.GetPotID(p.PUID)
	var query = ""
	var err error

	query = "SELECT nutrient_id, pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;"
	err = l.dbh.Get(&plant, query, potID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	plant.Nutrients = *l.GetPlantNutrients(plant.NutrientID)

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
		return nil, err
	}
	return jsonBytes, nil
}
