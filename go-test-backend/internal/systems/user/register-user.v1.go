package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	ID          string `json:"uuid"`
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

func (l *User) userAlreadyRegistered(id string) bool {
	var query = "SELECT COUNT(*) FROM users WHERE uuid=?;"
	_, err := l.dbh.Query(query, id)
	return err == sql.ErrNoRows
}

func (l *User) displayNameAlreadyUsed(name string) bool {
	var query = "SELECT COUNT(*) FROM users WHERE display_name=?;"
	_, err := l.dbh.Query(query, name)
	return err == sql.ErrNoRows
}

//RegisterUserV1 creates a register user object with given arguments
func (l *User) RegisterUserV1(ctx context.Context, p *RegisterUserV1Params, res *RegisterUserV1Result) error {
	if l.userAlreadyRegistered(p.ID) {
		return errors.New("user already registered")
	}
	if l.displayNameAlreadyUsed(p.DisplayName) {
		return errors.New("display name already used")
	}

	encrypted := l.encryptPassword(p.Password)

	var query = "INSERT INTO users (uuid, display_name, first_name, last_name, email, birthday, password_hash, registered_ts) VALUES (?,?,?,?,?,?,?,?);"
	_, err := l.dbh.Exec(query, p.ID, p.DisplayName, p.FirstName, p.LastName, p.Email, p.Birthday, encrypted, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
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
