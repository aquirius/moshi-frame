package user

import (
	"context"
	"time"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	redis "github.com/go-redis/redis/v8"
)

type DBUser struct {
	UUID          uint64 `db:"uuid"`
	Password      string `db:"password"`
	Authenticated bool
}

type SessionUser struct {
	UUID          string `db:"uuid"`
	Password      string `db:"password_hash"`
	Authenticated bool
}

// User
type LoginUser struct {
	ID string `json:"uuid"`
}

type LoginUserV1Params struct {
	Username  string  `json:"display_name"`
	Password  string  `json:"password"`
	SessionID *string `json:"session_id"`
}

type LoginUserV1Result struct {
	UUID      string `json:"uuid"`
	SessionID string `json:"session_id"`
}

func (c SessionUser) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

//LoginUserV1 verifys user login
func (l *User) LoginUserV1(ctx context.Context, p *LoginUserV1Params) (*LoginUserV1Result, error) {
	//res := &LoginUserV1Result{}
	var sessionID string
	var redisSession string
	var err error

	//generate random id
	sessionID = l.generateSessionID()

	if p.SessionID != nil {
		redisSession, err = l.rdb.Get(ctx, *p.SessionID).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
		//overwrite generated id if we have already one from request
		sessionID = *p.SessionID
	}
	ctx.Done()

	//set our redis context
	ctx = context.WithValue(ctx, sessionID, redisSession)

	var session SessionUser
	//get password hash from sql and check params
	err = l.dbh.Get(&session, "SELECT uuid, password_hash FROM users WHERE display_name=?;", p.Username)
	if err != nil {
		panic(err)
	}

	//check password encryption
	encrypted := l.encryptPassword(p.Password)

	if encrypted != session.Password {
		return nil, errors.New("passwords do not match")
	}

	//we authenticate the user if passwords do match
	session.Authenticated = true

	//connect user with newly generated session id
	err = l.rdb.Set(ctx, sessionID, session, 180*time.Second).Err()
	if err != nil {
		panic(err)
	}
	return &LoginUserV1Result{UUID: session.UUID, SessionID: sessionID}, nil
}

//LoginUserHandler handles login user request
func (l *User) LoginUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &LoginUserV1Params{}

	cookie, _ := r.Cookie("session-id")

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	ctx := context.Background()

	//if we have a session id store it to req body
	if cookie != nil && cookie.Value != "" {
		req.SessionID = &cookie.Value
	}

	session, err := l.LoginUserV1(ctx, req)
	if err != nil {
		return nil, err
	}

	//if we get a session id from login we store it as cookie
	addSessionCookie := &http.Cookie{
		Name:   "session-id",
		Value:  session.SessionID,
		MaxAge: 180,
	}

	http.SetCookie(w, addSessionCookie)

	//if we get a session id from login we store it as cookie
	addUserIDCookie := &http.Cookie{
		Name:   "uuid",
		Value:  session.UUID,
		MaxAge: 180,
	}

	http.SetCookie(w, addUserIDCookie)

	jsonBytes, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
