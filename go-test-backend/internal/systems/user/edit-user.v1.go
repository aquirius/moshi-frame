package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EditUser
type EditUser struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Title           string `json:"title"`
	Salutation      string `json:"salutation"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DisplayLanguage string `json:"language"`
	Country         string `json:"country"`
	Password        string `json:"password"`
}

// EditUserV1Params
type EditUserV1Params struct {
	ID          string  `json:"uuid"`
	Email       *string `json:"email"`
	DisplayName *string `json:"display_name"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Password    *string `json:"password"`
}

// EditUserV1 edits a user with given arguments
func (l *Users) EditUserV1(p *EditUserV1Params) error {
	var query string
	var arguments = []any{}
	fmt.Println(*p.DisplayName, *p.FirstName, *p.LastName, *p.Email)
	if *p.DisplayName == "" && *p.FirstName == "" && *p.LastName == "" && *p.Email == "" {
		return nil
	}
	query += "UPDATE users SET "
	if *p.DisplayName != "" {
		arguments = append(arguments, *p.DisplayName)
		query += "display_name=?,"
	}
	if *p.FirstName != "" {
		arguments = append(arguments, *p.FirstName)
		query += "first_name=?,"
	}
	if *p.LastName != "" {
		arguments = append(arguments, *p.LastName)
		query += "last_name=?,"
	}
	if *p.Email != "" {
		arguments = append(arguments, *p.Email)
		query += "email=?,"
	}

	//delete comma
	query = query[:len(query)-1]
	query += " WHERE uuid=?;"

	//uuid as last argument
	arguments = append(arguments, p.ID)

	_, err := l.dbh.Exec(query, arguments...)
	if err != nil {
		return err
	}
	return nil
}

// EditUserHandler handles editing one user
func (l *Users) EditUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &EditUserV1Params{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)

	err := l.EditUserV1(req)
	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal("success")
	fmt.Println(err)

	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
