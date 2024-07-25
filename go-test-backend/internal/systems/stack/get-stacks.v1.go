package stack

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-backend/m/v2/internal/systems/greenhouse"

	"github.com/gorilla/mux"
)

// GetUser
type GetStacks struct {
	SUID uint64 `db:"suid"`
}

// GetUserV1Params
type GetStacksV1Params struct {
	UUID uint64 `json:"uuid"`
	GUID uint64 `json:"guid"`
}

// GetUserV1Result
type GetStacksV1Result struct {
	Stacks []GetStacks `json:"stacks"`
}

// GetUserV1 gets user by uuid
func (l *Stack) GetStacksV1(ctx context.Context, p *GetStacksV1Params) (*GetStacksV1Result, error) {
	stacks := []uint64{}
	//_ = ctx.Value("stacks_id")
	greenhouse := greenhouse.NewGreenhouseProvider(ctx, l.dbh, l.rdb, "")
	greenhouseID := greenhouse.Greenhouse.GetGreenhouseID(p.GUID)

	err := l.dbh.Select(&stacks, "SELECT suid FROM stacks WHERE greenhouse_id=?;", greenhouseID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	getStacks := []GetStacks{}
	for _, v := range stacks {
		res := GetStacks{
			SUID: v,
		}
		getStacks = append(getStacks, res)
	}

	return &GetStacksV1Result{Stacks: getStacks}, nil
}

// GetUserHandler handles get user request
func (l *Stack) GetStacksHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)

	req := &GetStacksV1Params{
		UUID: uuid,
		GUID: guid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetStacksV1(ctx, req)
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
