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

// DeleteNotification nuid, created_ts, checked_ts, done_ts, title
type DeleteNotification struct {
	NUID      uint64 `db:"nuid"`
	CreatedTS int    `db:"created_ts"`
	CheckedTS uint64 `db:"checked_ts"`
	DoneTS    uint64 `db:"done_ts"`
	Title     string `db:"title"`
	Message   string `db:"message"`
}

// DeleteNotificationV1Params
type DeleteNotificationV1Params struct {
	UUID uint64 `json:"uuid"`
	NUID uint64 `json:"nuid"`
}

// DeleteNotificationV1Result
type DeleteNotificationV1Result struct {
}

// DeleteNotificationV1 gets pots by suid
func (l *Notification) DeleteNotificationV1(ctx context.Context, p *DeleteNotificationV1Params) (*DeleteNotificationV1Result, error) {
	var query = ""
	var err error

	query = "DELETE FROM notifications WHERE nuid=?;"
	_, err = l.dbh.Exec(query, p.NUID)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return nil, nil
}

// DeleteNotificationHandler handles get user request
func (l *Notification) DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	vars := mux.Vars(r)
	ctx := context.Background()

	uuid, _ := strconv.ParseUint(vars["uuid"], 0, 64)
	nuid, _ := strconv.ParseUint(vars["nuid"], 0, 64)

	req := &DeleteNotificationV1Params{
		UUID: uuid,
		NUID: nuid,
	}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	res, err := l.DeleteNotificationV1(ctx, req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
