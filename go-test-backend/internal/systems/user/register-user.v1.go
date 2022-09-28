package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// RegisterUser
type RegisterUser struct {
	ID          uint64 `json:"uuid"`
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Birthday    uint64 `json:"birthday"`
	Password    string `json:"password"`
}

//RegisterUserV1Params
type RegisterUserV1Params struct {
	ID          uint64 `json:"uuid"`
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Birthday    uint64 `json:"birthday"`
	Password    string `json:"password"`
}

//RegisterUserV1Result
type RegisterUserV1Result struct {
	User *RegisterUser `json:"user"`
}

func (l *User) existingUUID(uuid uint64) bool {
	var query = "SELECT id FROM users WHERE uuid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (l *User) existingUsername(name string) bool {
	var query = "SELECT id FROM users WHERE display_name=?;"
	var id int
	err := l.dbh.Get(&id, query, name)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

//RegisterUserV1 creates a register user object with given arguments
func (l *User) RegisterUserV1(ctx context.Context, p *RegisterUserV1Params, res *RegisterUserV1Result) error {
	if l.existingUUID(p.ID) {
		return errors.New("user already registered")
	}

	if l.existingUsername(p.DisplayName) {
		return errors.New("display name already taken")
	}

	encrypted := l.encryptPassword(p.Password)

	var query = "INSERT INTO users (uuid, display_name, first_name, last_name, email, birthday, password_hash, registered_ts) VALUES (?,?,?,?,?,?,?,?);"
	_, err := l.dbh.Exec(query, p.ID, p.DisplayName, p.FirstName, p.LastName, p.Email, p.Birthday, encrypted, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

//RegisterUserHandler handles register user request
func (l *User) RegisterUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	res := &RegisterUserV1Result{}
	req := &RegisterUserV1Params{}

	//read body and map on params
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)

	ctx := context.Background()

	err := l.RegisterUserV1(ctx, req, res)
	if err != nil {
		return nil, err
	}
	//build response body
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
