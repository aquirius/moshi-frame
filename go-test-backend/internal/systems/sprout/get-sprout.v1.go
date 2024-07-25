package sprout

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"test-backend/m/v2/internal/systems/stack"
)

// GetUser
type GetSprout struct {
	SproutUID uint64  `db:"sproutuid"`
	PH        float64 `db:"pH"`
	TDS       int64   `db:"TDS"`
	ORP       float64 `db:"ORP"`
	WaterTemp float64 `db:"h2oTemp"`
	AirTemp   float64 `db:"airTemp"`
	Humidity  float64 `db:"humidity"`
}

// GetUserV1Params
type GetSproutV1Params struct {
	SUID uint64 `json:"suid"`
}

// GetUserV1Result
type GetSproutV1Result struct {
	Sprout GetSprout `json:"sprout"`
}

// GetUserV1 gets user by uuid
func (l *Sprout) GetSproutV1(ctx context.Context, p *GetSproutV1Params) (*GetSproutV1Result, error) {
	sprout := GetSprout{}
	stack := stack.NewStackProvider(ctx, l.dbh, l.rdb, "")
	stackID := stack.Stack.GetStackID(p.SUID)

	fmt.Println(stackID)

	err := l.dbh.Select(&sprout, "SELECT sproutuid, pH, TDS, ORP, h2oTemp, airTemp, humidity FROM sprouts WHERE stack_id=?;", stackID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	fmt.Println(sprout)

	return &GetSproutV1Result{Sprout: sprout}, nil
}

// GetUserHandler handles get user request
func (l *Sprout) GetSproutHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	ctx := context.Background()

	req := &GetSproutV1Params{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	fmt.Println(req)
	res, err := l.GetSproutV1(ctx, req)
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