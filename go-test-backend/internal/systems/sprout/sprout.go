package sprout

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// SproutProvider provides *Sprout
type SproutProvider struct {
	Sprout *Sprout
}

// Sprout is capable of providing core access
type Sprout struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewSproutProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *SproutProvider {
	return &SproutProvider{
		&Sprout{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *SproutProvider) NewSprout() *Sprout {
	return b.Sprout
}

// serves user methods
func (sprout *Sprout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("post sprout get")
		res, err := sprout.GetSproutHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "add" {
			fmt.Println("post sprout add")
			res, err = sprout.AddSproutHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		if method == "get" {
			fmt.Println("post sprout get")
			res, err = sprout.GetSproutHandler(w, r)
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
