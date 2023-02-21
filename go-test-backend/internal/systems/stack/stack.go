package stack

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// StackProvider provides *Stack
type StackProvider struct {
	Stack *Stack
}

// Plant is capable of providing core access
type Stack struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewStackProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *StackProvider {
	return &StackProvider{
		&Stack{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *StackProvider) NewStack() *Stack {
	return b.Stack
}

func (l *Stack) GetStackID(suid uint64) int {
	var query = "SELECT id FROM stacks WHERE suid=?;"
	var id int
	err := l.dbh.Get(&id, query, suid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

//serves user methods
func (stack *Stack) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get stacks")
		res, err := stack.GetStacksHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
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
			fmt.Println("post stack add")
			res, err = stack.AddStackHandler(w, r)
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
