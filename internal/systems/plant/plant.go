package plant

import (
	"context"
	"database/sql"
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

// PlantProvider provides *Plant
type PlantsProvider struct {
	Plants *Plants
}

// Plant is capable of providing core access
type Plants struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewPlantsProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *PlantsProvider {
	return &PlantsProvider{
		&Plants{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *PlantProvider) NewPlant() *Plant {
	return b.Plant
}

func (b *PlantsProvider) NewPlants() *Plants {
	return b.Plants
}

func (l *Plant) GetPlantID(pluid uint64) int {
	var query = "SELECT id FROM plants WHERE pluid=?;"
	var id int
	err := l.dbh.Get(&id, query, pluid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

func (l *Plant) GetPotID(puid uint64) int {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := l.dbh.Get(&id, query, puid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

func (l *Plants) GetPotID(puid uint64) int {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := l.dbh.Get(&id, query, puid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

// serves user methods
func (plant *Plant) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "add" {
			res, err = plant.AddPlantHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		if method == "harvest" {
			res, err = plant.HarvestPlantHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodGet:
		res, err := plant.GetPlantHandler(w, r)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusOK)
				w.Write(res)
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
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

func (plants *Plants) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "get-many" {
			res, err = plants.GetPlantsHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return

	default:
		return
	}
}
