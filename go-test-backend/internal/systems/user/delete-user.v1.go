package user

import (
	"encoding/json"
	"io"
	"net/http"
)

// DeleteUser
type DeleteUser struct {
	ID string `json:"uuid"`
}

// DeleteUserV1Params
type DeleteUserV1Params struct {
	ID string `json:"uuid"`
}

// DeleteUsersV1Result
type DeleteUsersV1Result struct {
}

// DeleteUserV1 deletes a user with given uuid
func (l *Users) DeleteUserV1(p *DeleteUserV1Params) error {
	_, err := l.dbh.Exec("DELETE FROM users WHERE uuid=?;", p.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserHandler handles the deletion of a user
func (l *Users) DeleteUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &DeleteUserV1Params{}

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	err := l.DeleteUserV1(req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal("success")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
