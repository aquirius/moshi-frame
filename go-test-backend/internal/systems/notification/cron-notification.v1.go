package notification

import (
	"database/sql"
	"fmt"

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

// GetNotificationV1 gets pots by suid
func (l *Notifications) CronNotification() error {
	var err error

	c := cron.New()
	cronID, err := c.AddFunc("@every 1s", func() {
		var query = ""

		notification := Status{}

		query += "SELECT "
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

		//query = "SELECT nuid, created_ts, checked_ts, done_ts, title, message FROM notifications WHERE nuid=?;"
		err = l.dbh.Get(&notification, query, 1)
		fmt.Println(notification, err)

		if err == sql.ErrNoRows {
			return
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
