package user

import (
	"context"

	"net/http"
)

// User
type LogoutUser struct {
	ID uint64 `json:"uuid"`
}

type LogoutUserV1Params struct {
	UUID      string `json:"uuid"`
	SessionID string `json:"session_id"`
}

type LogoutUserV1Result struct {
}

// LogoutUserV1 verifys user logout
func (l *User) LogoutUserV1(ctx context.Context, p *LogoutUserV1Params) (*LogoutUserV1Result, error) {
	//we have no session id so we are already logged out
	if p.SessionID == "" {
		return nil, nil
	}

	l.rdb.Del(ctx, p.SessionID)

	return nil, nil
}

// LogoutUserHandler handles logout user request
func (l *User) LogoutUserHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &LogoutUserV1Params{}
	ctx := context.Background()

	delSessionCookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: 0,
	}
	http.SetCookie(w, delSessionCookie)

	delUserIDCookie := &http.Cookie{
		Name:   "uuid",
		Value:  "",
		MaxAge: 0,
	}
	http.SetCookie(w, delUserIDCookie)

	_, err := l.LogoutUserV1(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
