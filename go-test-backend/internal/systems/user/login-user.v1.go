package user

import (
	"context"
	"io"
	"time"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	redis "github.com/go-redis/redis/v8"
)

// User
type LoginUser struct {
	ID uint64 `json:"uuid"`
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

// LoginUserV1 verifys user login
func (l *User) LoginUserV1(ctx context.Context, p *LoginUserV1Params) (*LoginUserV1Result, error) {
	//res := &LoginUserV1Result{}
	var sessionID string
	var redisSession string
	var err error

	if !l.existingUsername(p.Username) {
		return nil, errors.New("display name does not exist")
	}

	if p.SessionID != nil {
		redisSession, err = l.rdb.Get(ctx, *p.SessionID).Result()
		if err == redis.Nil {
			fmt.Println("token does not exist")
		} else if err != nil {
			panic(err)
		}
		//overwrite generated id if we have already one from request
		sessionID = *p.SessionID
	} else {
		sessionID = l.generateSessionID()
	}

	//set our redis context
	ctx = context.WithValue(ctx, sessionID, redisSession)

	var session SessionUser
	//get password hash from sql and check params
	err = l.dbh.Get(&session, "SELECT id, uuid, password_hash FROM users WHERE display_name=?;", p.Username)
	if err != nil {
		return nil, errors.New("no user found")
	}

	//check password encryption
	encrypted := l.encryptPassword(p.Password)

	if encrypted != session.Password {
		return nil, errors.New("passwords do not match")
	}

	//we authenticate the user if passwords do match
	session.Authenticated = true

	//connect user with newly generated session id
	_, err = l.rdb.Set(ctx, sessionID, session, 24*60*60*time.Second).Result()
	if err != nil {
		return nil, errors.New("redis err saving session")
	}

	ctx = context.WithValue(ctx, "user_id", session.ID)

	return &LoginUserV1Result{UUID: session.UUID, SessionID: sessionID}, nil
}

// LoginUserHandler handles login user request
func (l *User) LoginUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &LoginUserV1Params{}

	cookie, _ := r.Cookie("session-id")

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)

	ctx := context.Background()
	//if we have a session id store it to req body
	if req.SessionID == nil && (cookie != nil && cookie.Value != "") {
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
		MaxAge: 24 * 60 * 60,
	}
	http.SetCookie(w, addSessionCookie)

	//if we get a user id from login we store it as cookie
	addUserIDCookie := &http.Cookie{
		Name:   "uuid",
		Value:  session.UUID,
		MaxAge: 24 * 60 * 60,
	}
	http.SetCookie(w, addUserIDCookie)

	jsonBytes, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
