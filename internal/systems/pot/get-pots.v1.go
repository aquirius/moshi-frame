package pot

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"test-backend/m/v2/internal/systems/stack"
)

// GetPots
type GetPots struct {
	PUID uint64 `db:"puid"`
}

// GetPotsV1Params
type GetPotsV1Params struct {
	SUID uint64 `json:"suid"`
}

// GetPotsV1Result
type GetPotsV1Result struct {
	Pots []GetPots `json:"pots"`
}

// GetPotsV1 gets pots by suid
func (l *Pot) GetPotsV1(ctx context.Context, p *GetPotsV1Params) (*GetPotsV1Result, error) {
	pots := []uint64{}
	stack := stack.NewStackProvider(ctx, l.dbh, l.rdb, "")
	stackID := stack.Stack.GetStackID(p.SUID)

	err := l.dbh.Select(&pots, "SELECT puid FROM pots WHERE stack_id=?;", stackID)
	if err == sql.ErrNoRows {
		return nil, err
	}
	getPots := []GetPots{}
	for _, v := range pots {
		res := GetPots{
			PUID: v,
		}
		getPots = append(getPots, res)
	}

	return &GetPotsV1Result{Pots: getPots}, nil
}

// GetUserHandler handles get user request
func (l *Pot) GetPotsHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	ctx := context.Background()

	req := &GetPotsV1Params{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetPotsV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Fatal("error in json")
		return nil, err
	}
	return jsonBytes, nil
}
