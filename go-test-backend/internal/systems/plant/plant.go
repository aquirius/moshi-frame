package plant

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// PlantProvider provides *Plant
type PlantProvider struct {
	Plant *Plant
}

// Plant is capable of providing core access
type Plant struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewPlantProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *PlantProvider {
	return &PlantProvider{
		&Plant{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *PlantProvider) NewPlant() *Plant {
	return b.Plant
}

func (l *Plant) GetPlantID(pluid uint64) int {
	var query = "SELECT id FROM greenhouse WHERE guid=?;"
	var id int
	err := l.dbh.Get(&id, query, pluid)
	fmt.Println("getplant", pluid, err, id)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

//serves user methods
func (plant *Plant) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		fmt.Println("add plant")
		res, err := plant.AddPlantHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	default:
		return
	}
}
