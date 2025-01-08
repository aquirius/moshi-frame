package sprout

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/stack"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetUser
type AddSprout struct {
	SPUID uint64 `json:"spuid"`
}

// GetUserV1Params
type AddSproutV1Params struct {
	SUID uint64 `json:"suid"`
}

// GetUserV1Result
type AddSproutV1Result struct {
	Sprout AddSprout `json:"sprout"`
}

// GetUserV1 gets user by uuid
func (l *Sprout) AddSproutV1(ctx context.Context, p *AddSproutV1Params) (*AddSproutV1Result, error) {
	sproutuid, err := uuid.NewRandom()
	stack := stack.NewStackProvider(ctx, l.dbh, l.rdb, "")
	stackID := stack.Stack.GetStackID(p.SUID)

	query := "INSERT INTO sprouts (sproutuid, stack_id, pH, TDS, ORP, h2oTemp, airTemp, humidity) VALUES (?,?,?,?,?,?,?,?);"
	_, err = l.dbh.Exec(query, sproutuid.ID(), stackID, 5.5, 1000, 450, 20.3, 26.5, 69)
	if err != nil {
		return nil, err
	}

	res := &AddSprout{
		SPUID: uint64(sproutuid.ID()),
	}

	return &AddSproutV1Result{Sprout: *res}, nil
}

// GetUserHandler handles get user request
func (l *Sprout) AddSproutHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	suid, _ := strconv.ParseUint(vars["suid"], 0, 32)

	req := &AddSproutV1Params{
		SUID: suid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.AddSproutV1(ctx, req)
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
