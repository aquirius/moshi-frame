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

type Crop struct {
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

// GetPlant
type GetPlant struct {
	CropID      int    `db:"crop_id"`
	NutrientID  int    `db:"nutrient_id"`
	PLUID       uint64 `db:"pluid"`
	CreatedTS   uint64 `db:"created_ts"`
	PlantedTS   uint64 `db:"planted_ts"`
	HarvestedTS uint64 `db:"harvested_ts"`

	Nutrients Nutrients
	Crop      Crop
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
	return &nutrients
}

func (l *Plant) GetPlantCrop(cropID int) *Crop {
	var query = "SELECT cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max FROM crops WHERE id=?;"
	crop := Crop{}
	err := l.dbh.Get(&crop, query, cropID)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return &crop
}

// GetPlantV1 gets pots by suid
func (l *Plant) GetPlantV1(ctx context.Context, p *GetPlantV1Params) (*GetPlantV1Result, error) {
	plant := GetPlant{}

	pot := pot.NewPotProvider(ctx, l.dbh, l.rdb, "")
	potID := pot.Pot.GetPotID(p.PUID)
	var query = ""
	var err error

	query = "SELECT crop_id, nutrient_id, pluid, created_ts, planted_ts, harvested_ts FROM plants WHERE pot_id=?;"
	err = l.dbh.Get(&plant, query, potID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	plant.Nutrients = *l.GetPlantNutrients(plant.NutrientID)
	plant.Crop = *l.GetPlantCrop(plant.CropID)

	fmt.Println(plant.Crop)

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
