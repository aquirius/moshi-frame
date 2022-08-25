package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//GetUserV1Params
type GetUserV1Params struct {
	UUID string `json:"uuid"`
}

//GetUserV1Result
type GetUserV1Result struct {
	User GetUser `json:"user"`
}

//GetUserV1 gets user by uuid
func (l *User) GetUserV1(ctx context.Context, p *GetUserV1Params) (*GetUserV1Result, error) {
	user := GetUser{}
	v := ctx.Value("session-id")
	fmt.Println(v)
	err := l.dbh.Get(&user, "SELECT uuid, registered_ts, display_name, first_name, last_name, email, birthday FROM users WHERE uuid=?", p.UUID)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	return &GetUserV1Result{User: user}, nil
}

//GetUserHandler handles get user request
func (l *User) GetUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	cookie, _ := r.Cookie("session-id")
	ctx := context.Background()

	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		ctx = context.WithValue(ctx, "session-id", cookie.Value)
	}

	req := &GetUserV1Params{
		UUID: vars["uuid"],
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetUserV1(ctx, req)
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
