package greenhouse

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// GreenhouseProvider provides *Greenhouse
type GreenhouseProvider struct {
	Greenhouse *Greenhouse
}

// Greenhouse is capable of providing core access
type Greenhouse struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// GreenhouseProvider provides *Greenhouse
type GreenhousesProvider struct {
	Greenhouses *Greenhouses
}

// Greenhouse is capable of providing core access
type Greenhouses struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

func (l *Greenhouses) ExistingPUID(uuid uint32) bool {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (l *Greenhouse) GetGreenhouseID(guid uint64) int {
	var query = "SELECT id FROM greenhouses WHERE guid=?;"
	var id int
	err := l.dbh.Get(&id, query, guid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

// NewCoreProvider returns a new Core provider
func NewGreenhouseProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *GreenhouseProvider {
	return &GreenhouseProvider{
		&Greenhouse{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

// NewCoreProvider returns a new Core provider
func NewGreenhousesProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *GreenhousesProvider {
	return &GreenhousesProvider{
		&Greenhouses{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *GreenhouseProvider) NewGreenhouse() *Greenhouse {
	return b.Greenhouse
}

// NewUsers
func (b *GreenhousesProvider) NewGreenhouses() *Greenhouses {
	return b.Greenhouses
}

// serves user methods
func (b *Greenhouse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get greenhouse")
		res, err := b.GetGreenhouseHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodPut:
		fmt.Println("put greenhouse")
		res, err := b.EditGreenhouseHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	default:
		return
	}
}

// serves Greenhouses methods
func (b *Greenhouses) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get greenhouses")
		res, err := b.GetGreenhousesHandler(w, r)
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
			fmt.Println("post Greenhouse add")
			res, err = b.AddGreenhouseHandler(w, r)
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
