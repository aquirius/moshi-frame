package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//GetUser
type GetUser struct {
	UUID        string `db:"uuid"`
	TS          string `db:"registered_ts"`
	DisplayName string `db:"display_name"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Birthday    string `db:"birthday"`
}

//GetUsersV1Params
type GetUsersV1Params struct {
}

//GetUsersV1Result
type GetUsersV1Result struct {
	Users []GetUser `json:"users"`
}

//GetUsersV1 gets all users
func (l *Users) GetUsersV1(p *GetUsersV1Params) (*GetUsersV1Result, error) {
	users := []GetUser{}

	err := l.dbh.Select(&users, "SELECT uuid, registered_ts, display_name, first_name, last_name, email, birthday FROM users;")
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil, err
	}

	return &GetUsersV1Result{Users: users}, nil
}

//GetUsersHandler get all users request
func (l *Users) GetUsersHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &GetUsersV1Params{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetUsersV1(nil)
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
