package crop

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// CropProvider provides *Crop
type CropProvider struct {
	Crop *Crop
}

// Crop is capable of providing core access
type Crop struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewCropProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *CropProvider {
	return &CropProvider{
		&Crop{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

// CropProvider provides *Crop
type CropsProvider struct {
	Crops *Crops
}

// Crop is capable of providing core access
type Crops struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewCropsProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *CropsProvider {
	return &CropsProvider{
		&Crops{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *CropProvider) NewCrop() *Crop {
	return b.Crop
}

func (b *CropsProvider) NewCrops() *Crops {
	return b.Crops
}

func (l *Crop) GetCropID(pluid uint64) int {
	var query = "SELECT id FROM crops WHERE pluid=?;"
	var id int
	err := l.dbh.Get(&id, query, pluid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

// serves user methods
func (crop *Crop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "add" {
			res, err = crop.AddCropHandler(w, r)
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
		res, err := crop.GetCropHandler(w, r)
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

func (crops *Crops) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "get-many" {
			res, err = crops.GetCropsHandler(w, r)
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
