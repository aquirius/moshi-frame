package plant

import (
	"context"
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

//serves user methods
func (plant *Plant) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get plants")
		res, err := plant.GetGreenhouseHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "add" {
			fmt.Println("post plant add")
			res, err = plant.AddGreenhouseHandler(w, r)
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
