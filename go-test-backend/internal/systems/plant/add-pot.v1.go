package plant

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//GetUser
type AddPot struct {
	PUID uint64 `json:"puid"`
}

//GetUserV1Params
type AddPotV1Params struct {
	UUID uint64 `json:"uuid"`
	SUID uint64 `json:"suid"`
}

//GetUserV1Result
type AddPotV1Result struct {
	Pot AddPot `json:"pot"`
}

//GetUserV1 gets user by uuid
func (l *Plant) AddPotV1(ctx context.Context, p *AddPotV1Params) (*AddPotV1Result, error) {
	puid, err := uuid.NewUUID()
	userID := l.getUserID(p.UUID)
	stackID := l.getStackID(p.SUID)

	fmt.Println("add pot params : ", p)

	query := "INSERT INTO pots (puid, stack_id, user_id) VALUES (?,?,?);"
	_, err = l.dbh.Exec(query, puid.ID(), stackID, userID)
	if err != nil {
		return nil, err
	}

	res := &AddPot{
		PUID: uint64(puid.ID()),
	}

	return &AddPotV1Result{Pot: *res}, nil
}

//GetUserHandler handles get user request
func (l *Plant) AddPotHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 32)
	req := &AddPotV1Params{
		UUID: uuid,
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.AddPotV1(ctx, req)
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
