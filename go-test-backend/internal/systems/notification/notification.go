package notification

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// NotificationProvider provides *Notification
type NotificationProvider struct {
	Notification *Notification
}

// Notification is capable of providing core access
type Notification struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewNotificationProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *NotificationProvider {
	return &NotificationProvider{
		&Notification{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

// NotificationProvider provides *Notification
type NotificationsProvider struct {
	Notifications *Notifications
}

// Notification is capable of providing core access
type Notifications struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// NewCoreProvider returns a new Core provider
func NewNotificationsProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *NotificationsProvider {
	return &NotificationsProvider{
		&Notifications{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

func (b *NotificationProvider) NewNotification() *Notification {
	return b.Notification
}

func (b *NotificationsProvider) NewNotifications() *Notifications {

	b.Notifications.CronNotification()

	return b.Notifications
}

func (l *Notification) GetNotificationID(nuid uint64) int {
	var query = "SELECT id FROM notifications WHERE nuid=?;"
	var id int
	err := l.dbh.Get(&id, query, nuid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

// serves user methods
func (notification *Notification) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get plant")
		res, err := notification.GetNotificationHandler(w, r)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusOK)
				w.Write(res)
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return

	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		fmt.Println(method)
		if method == "delete" {
			res, err = notification.DeleteNotificationHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	default:
		return
	}
}

func (plants *Notifications) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		fmt.Println(method)
		if method == "get-many" {
			res, err = plants.GetNotificationsHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return

	default:
		return
	}
}
