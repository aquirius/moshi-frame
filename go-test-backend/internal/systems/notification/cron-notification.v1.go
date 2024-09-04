package notification

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
)

/*
SELECT
    customers.name AS customer_name,
    orders.order_id,
    orders.order_date,
    products.product_name,
    products.price,
    suppliers.supplier_name
FROM
    customers
JOIN
    orders ON customers.customer_id = orders.customer_id
JOIN
    products ON orders.product_id = products.product_id
JOIN
    suppliers ON products.product_id = suppliers.product_id;
*/

type Status struct {
	GUID    uint64 `db:"guid"`
	Address string `db:"address"`
	Zip     uint64 `db:"zip"`

	DisplayName string  `db:"display_name"`
	Status      string  `db:"status"`
	Destination string  `db:"destination"`
	TempIn      float64 `db:"tempIn"`
	TempOut     float64 `db:"tempOut"`
	Humidity    float64 `db:"humidity"`
	Brightness  float64 `db:"brightness"`
	Co2         float64 `db:"co2"`

	SUID uint64 `db:"suid"`

	SproutUID      uint64  `db:"sproutuid"`
	PH             float64 `db:"pH"`
	TDS            int64   `db:"TDS"`
	ORP            float64 `db:"ORP"`
	WaterTemp      float64 `db:"h2oTemp"`
	AirTemp        float64 `db:"airTemp"`
	HumiditySprout float64 `db:"humidity"`

	PUID uint64 `db:"puid"`

	CropID      int    `db:"crop_id"`
	NutrientID  int    `db:"nutrient_id"`
	PLUID       uint64 `db:"pluid"`
	CreatedTS   uint64 `db:"created_ts"`
	PlantedTS   uint64 `db:"planted_ts"`
	HarvestedTS uint64 `db:"harvested_ts"`

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

type Crop struct {
	CropID       int     `db:"id"`
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

type CropsRange struct {
	CropName     float64 `db:"crop_name"`
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

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func generate10DigitNumber() int64 {
	rand.Seed(time.Now().UnixNano())

	min := int64(1_000_000_000)
	max := int64(9_999_999_999)

	return rand.Int63n(max-min+1) + min
}

func (l *Notifications) InsertNotificationMessage(userID int, message string) error {
	nuid := generate10DigitNumber()
	var query = ""
	query += "INSERT INTO notifications "
	query += "(nuid, created_ts, checked_ts, done_ts, title, message, user_id) "
	query += "VALUES(?,?,?,?,?,?,?);"
	_, err := l.dbh.Exec(query, nuid, time.Now().Unix(), time.Now().Unix(), time.Now().Unix(), "Notification", message, userID)
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func (l *Notifications) GetAllByUser(uuid uint64) []Status {
	var query = ""
	status := []Status{}

	query += "SELECT DISTINCT "
	query += "greenhouses.guid, greenhouses.address, greenhouses.zip, greenhouses.display_name, greenhouses.status, greenhouses.destination, greenhouses.tempIn, greenhouses.tempOut, greenhouses.humidity, greenhouses.brightness, greenhouses.co2, "
	query += "stacks.suid, "
	query += "sprouts.sproutuid, sprouts.pH, sprouts.TDS, sprouts.ORP, sprouts.h2oTemp, sprouts.airTemp, sprouts.humidity, "
	query += "pots.puid, "
	query += "plants.crop_id, plants.nutrient_id, plants.pluid, plants.created_ts, plants.planted_ts, plants.harvested_ts, "
	query += "crops.cuid, crops.crop_name, crops.air_temp_min, crops.air_temp_max, crops.humidity_min, crops.humidity_max, crops.ph_level_min, crops.ph_level_max, crops.orp_min, crops.orp_max, crops.tds_min, crops.tds_max, crops.water_temp_min, crops.water_temp_max "
	query += "FROM users_greenhouses "
	query += "JOIN greenhouses ON greenhouses.id = users_greenhouses.greenhouse_id "
	query += "JOIN stacks ON greenhouses.id = stacks.greenhouse_id "
	query += "JOIN sprouts ON stacks.id = sprouts.stack_id "
	query += "JOIN pots ON stacks.id = pots.stack_id "
	query += "JOIN plants ON pots.id = plants.pot_id "
	query += "JOIN crops ON crops.id = plants.crop_id "
	query += "WHERE users_greenhouses.user_id=?; "

	/*
		SELECT greenhouses.guid, greenhouses.address, greenhouses.zip, greenhouses.display_name, greenhouses.status, greenhouses.destination, greenhouses.tempIn, greenhouses.tempOut, greenhouses.humidity, greenhouses.brightness, greenhouses.co2, stacks.suid, sprouts.sproutuid, sprouts.pH, sprouts.TDS, sprouts.ORP, sprouts.h2oTemp, sprouts.airTemp, sprouts.humidity, pots.puid, plants.crop_id, plants.nutrient_id, plants.pluid, plants.created_ts, plants.planted_ts, plants.harvested_ts, crops.cuid, crops.crop_name, crops.air_temp_min, crops.air_temp_max, crops.humidity_min, crops.humidity_max, crops.ph_level_min, crops.ph_level_max, crops.orp_min, crops.orp_max, crops.tds_min, crops.tds_max, crops.water_temp_min, crops.water_temp_max FROM greenhouses JOIN stacks ON greenhouses.id = stacks.greenhouse_id JOIN sprouts ON stacks.id = sprouts.stack_id JOIN pots ON stacks.id = pots.stack_id JOIN plants ON pots.id = plants.pot_id JOIN crops ON crops.id = plants.crop_id;
	*/
	err := l.dbh.Select(&status, query, uuid)
	if err == sql.ErrNoRows {
		return nil
	}
	return status
}

func (l *Notifications) GetCropsSproutRange() *CropsRange {
	var query = ""
	crops := []Crop{}

	query += "SELECT "
	query += "id, cuid, crop_name, air_temp_min, air_temp_max, humidity_min, humidity_max, ph_level_min, ph_level_max, orp_min, orp_max, tds_min, tds_max, water_temp_min, water_temp_max "
	query += "FROM crops; "
	err := l.dbh.Select(&crops, query)
	if err == sql.ErrNoRows {
		return nil
	}
	cropsRange := CropsRange{
		AirTempMin:   min(crops[0].AirTempMin, crops[1].AirTempMin),
		AirTempMax:   max(crops[0].AirTempMax, crops[1].AirTempMax),
		HumidityMin:  min(crops[0].HumidityMin, crops[1].HumidityMin),
		HumidityMax:  max(crops[0].HumidityMax, crops[1].HumidityMax),
		PHLevelMin:   min(crops[0].PHLevelMin, crops[1].PHLevelMin),
		PHLevelMax:   max(crops[0].PHLevelMax, crops[1].PHLevelMax),
		OrpMin:       min(crops[0].OrpMin, crops[1].OrpMin),
		OrpMax:       min(crops[0].OrpMax, crops[1].OrpMax),
		TdsMin:       minInt(crops[0].TdsMin, crops[1].TdsMin),
		TdsMax:       maxInt(crops[0].TdsMax, crops[1].TdsMax),
		WaterTempMin: min(crops[0].WaterTempMin, crops[1].WaterTempMin),
		WaterTempMax: max(crops[0].WaterTempMax, crops[1].WaterTempMax),
	}
	return &cropsRange
}

//0  0    0 0 69 0 0 0 3262623532 5.5 1000 450 20.3 26.5 0 0 0 0 0 0 0 0 3454994911 lettuce 18 28 60 80 5 6 400 500 800 1200 18 22}
type MergedCropsPerSprout struct {
	SproutUID      map[uint64][]int
	PH             float64 `db:"pH"`
	TDS            int64   `db:"TDS"`
	ORP            float64 `db:"ORP"`
	WaterTemp      float64 `db:"h2oTemp"`
	AirTemp        float64 `db:"airTemp"`
	HumiditySprout float64 `db:"humidity"`
}

func groupIDsByValue(pairs []GetCropsPerSproutResult) []MergedCropsPerSprout {
	merged := []MergedCropsPerSprout{}
	result := make(map[uint64][]int)

	for _, pair := range pairs {
		result[pair.SproutUID] = append(result[pair.SproutUID], pair.CropID)
		merged = append(merged, MergedCropsPerSprout{
			SproutUID:      result,
			PH:             pair.PH,
			TDS:            pair.TDS,
			ORP:            pair.ORP,
			WaterTemp:      pair.WaterTemp,
			AirTemp:        pair.AirTemp,
			HumiditySprout: pair.HumiditySprout,
		})
	}

	return merged
}

type GetCropsPerSproutResult struct {
	CropID         int     `db:"crop_id"`
	SproutUID      uint64  `db:"sproutuid"`
	PH             float64 `db:"pH"`
	TDS            int64   `db:"TDS"`
	ORP            float64 `db:"ORP"`
	WaterTemp      float64 `db:"h2oTemp"`
	AirTemp        float64 `db:"airTemp"`
	HumiditySprout float64 `db:"humidity"`
}

func (l *Notifications) GetCropsPerSproutWithUser(userID int) []GetCropsPerSproutResult {
	var query = ""
	crops := []GetCropsPerSproutResult{}
	query += "SELECT DISTINCT "
	query += "plants.crop_id, sprouts.sproutuid, sprouts.pH, sprouts.TDS, sprouts.ORP, sprouts.h2oTemp, sprouts.airTemp, sprouts.humidity "
	query += "FROM users_greenhouses "
	query += "JOIN greenhouses ON greenhouses.id = users_greenhouses.greenhouse_id "
	query += "JOIN stacks ON greenhouses.id = stacks.greenhouse_id "
	query += "JOIN sprouts ON stacks.id = sprouts.stack_id "
	query += "JOIN pots ON stacks.id = pots.stack_id "
	query += "JOIN plants ON pots.id = plants.pot_id "
	query += "JOIN crops ON crops.id = plants.crop_id "
	query += "WHERE users_greenhouses.user_id=?; "
	err := l.dbh.Select(&crops, query, userID)
	if err == sql.ErrNoRows {
		return nil
	}
	return crops
}

// GetNotificationV1 gets pots by suid
func (l *Notifications) CronNotification() error {
	var err error
	c := cron.New()
	cronID, err := c.AddFunc("@every 10s", func() {
		var query = ""

		fmt.Println("hello")
		notifications := []Status{}

		query += "SELECT DISTINCT "
		query += "crops.cuid, crops.crop_name, crops.air_temp_min, crops.air_temp_max, crops.humidity_min, crops.humidity_max, crops.ph_level_min, crops.ph_level_max, crops.orp_min, crops.orp_max, crops.tds_min, crops.tds_max, crops.water_temp_min, crops.water_temp_max, "
		query += "sprouts.sproutuid, sprouts.pH, sprouts.TDS, sprouts.ORP, sprouts.h2oTemp, sprouts.airTemp, sprouts.humidity "
		query += "FROM users_greenhouses "
		query += "JOIN greenhouses ON greenhouses.id = users_greenhouses.greenhouse_id "
		query += "JOIN stacks ON greenhouses.id = stacks.greenhouse_id "
		query += "JOIN sprouts ON stacks.id = sprouts.stack_id "
		query += "JOIN pots ON stacks.id = pots.stack_id "
		query += "JOIN plants ON pots.id = plants.pot_id "
		query += "JOIN crops ON crops.id = plants.crop_id "
		query += "WHERE plants.planted_ts > 0 AND users_greenhouses.user_id=?; "

		err = l.dbh.Select(&notifications, query, 1)
		if err == sql.ErrNoRows {
			return
		}

		for _, notification := range notifications {
			msg := ""

			if notification.PH > notification.PHLevelMax {
				msg += "PH Level too high - Water fertilizer down"
			}
			if notification.PH < notification.PHLevelMin {
				msg += "PH Level too low - Fertilize water"
			}
			if notification.AirTemp > notification.AirTempMax {
				msg += "Air Temperature too high"
			}
			if notification.AirTemp < notification.AirTempMin {
				msg += "Air Temperature too low"
			}
			if notification.Humidity > notification.HumidityMax {
				msg += "Humidity too high"
			}
			if notification.Humidity < notification.HumidityMin {
				msg += "Humidity too low"
			}
			if notification.WaterTemp > notification.WaterTempMax {
				msg += "Water Temperature too high"
			}
			if notification.WaterTemp < notification.WaterTempMin {
				msg += "Water Temperature too low"
			}
			if msg != "" {
				l.InsertNotificationMessage(1, msg)
			}
		}
	})
	if err != nil {
		return err
	}

	c.Start()
	fmt.Println("initializing cron : ", cronID, c.Entries())

	//c.Stop()

	return nil
}
