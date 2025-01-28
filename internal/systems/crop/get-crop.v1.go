package crop

import (
	"context"
	"database/sql"
	"encoding/json"
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

// GetCrop
type GetCrop struct {
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

// GetCropV1Params
type GetCropV1Params struct {
	CUID uint64 `json:"cuid"`
}

// GetCropV1Result
type GetCropV1Result struct {
	Crop GetCrop `json:"crop"`
}

// GetCrop gets pots by suid
func (l *Crop) GetCropV1(ctx context.Context, p *GetCropV1Params) (*GetCropV1Result, error) {
	crop := GetCrop{}

	pot := pot.NewPotProvider(ctx, l.dbh, l.rdb, "")
	potID := pot.Pot.GetPotID(p.CUID)
	var query = ""
	var err error

	query = "SELECT crop_id, nutrient_id, pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;"
	err = l.dbh.Get(&crop, query, potID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &GetCropV1Result{Crop: crop}, nil
}

// GetCropHandler handles get crop request
func (l *Crop) GetCropHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	cuid, _ := strconv.ParseUint(vars["cuid"], 0, 32)
	req := &GetCropV1Params{
		CUID: cuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetCropV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
