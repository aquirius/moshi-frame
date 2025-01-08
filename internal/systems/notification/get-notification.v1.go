package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetNotification nuid, created_ts, checked_ts, done_ts, title
type GetNotification struct {
	NUID      uint64 `db:"nuid"`
	CreatedTS int    `db:"created_ts"`
	CheckedTS uint64 `db:"checked_ts"`
	DoneTS    uint64 `db:"done_ts"`
	Title     string `db:"title"`
	Message   string `db:"message"`
}

// GetNotificationV1Params
type GetNotificationV1Params struct {
	NUID uint64 `json:"nuid"`
}

// GetNotificationV1Result
type GetNotificationV1Result struct {
	Notification GetNotification `json:"notification"`
}

// GetNotificationV1 gets pots by suid
func (l *Notification) GetNotificationV1(ctx context.Context, p *GetNotificationV1Params) (*GetNotificationV1Result, error) {
	notification := GetNotification{}
	var query = ""
	var err error

	query = "SELECT nuid, created_ts, checked_ts, done_ts, title, message FROM notifications WHERE nuid=?;"
	err = l.dbh.Get(&notification, query, p.NUID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &GetNotificationV1Result{Notification: notification}, nil
}

// GetUserHandler handles get user request
func (l *Notification) GetNotificationHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	nuid, _ := strconv.ParseUint(vars["nuid"], 0, 32)
	req := &GetNotificationV1Params{
		NUID: nuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.GetNotificationV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
