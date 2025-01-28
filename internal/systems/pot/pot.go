package pot

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// PotProvider provides *Pot
type PotProvider struct {
	Pot *Pot
}

// Pot is capable of providing core access
type Pot struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewPotProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *PotProvider {
	return &PotProvider{
		&Pot{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *PotProvider) NewPot() *Pot {
	return b.Pot
}

func (pot *Pot) GetPotID(puid uint64) int {
	var query = "SELECT id FROM pots WHERE puid=?;"
	var id int
	err := pot.dbh.Get(&id, query, puid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

// serves user methods
func (pot *Pot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		return
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "add" {
			res, err = pot.AddPotHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		if method == "get" {
			res, err = pot.GetPotsHandler(w, r)
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
