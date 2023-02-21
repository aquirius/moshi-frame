package user

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// UserProvider provides *User
type UserProvider struct {
	User *User
}

// UsersProvider provides *Users
type UsersProvider struct {
	Users *Users
}

// User is capable of providing user access
type User struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

// Users is capable of providing users access
type Users struct {
	dbh *sqlx.DB
	rdb *redis.Client
}

type DBUser struct {
	UUID          uint64 `db:"uuid"`
	Password      string `db:"password"`
	Authenticated bool
}

type SessionUser struct {
	ID            uint64 `db:"id"`
	UUID          string `db:"uuid"`
	Password      string `db:"password_hash"`
	Authenticated bool
}

// NewUserProvider returns a new User provider
func NewUserProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *UserProvider {
	return &UserProvider{
		&User{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

// NewUserProvider returns a new Users provider
func NewUsersProvider(ctx context.Context, dbh *sqlx.DB, rdb *redis.Client, urlPrefixBackend string) *UsersProvider {
	return &UsersProvider{
		&Users{
			dbh: dbh,
			rdb: rdb,
		},
	}
}

//NewUser
func (b *UserProvider) NewUser() *User {
	return b.User
}

//NewUsers
func (b *UsersProvider) NewUsers() *Users {
	return b.Users
}

func (c *User) generateSessionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (c *User) encryptPassword(pw string) string {
	crypt := md5.New()
	io.WriteString(crypt, pw)
	return fmt.Sprintf("%x", crypt.Sum(nil))
}

func (c SessionUser) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (l *User) GetUserID(uuid uint64) int {
	var query = "SELECT id FROM users WHERE uuid=?;"
	var id int
	err := l.dbh.Get(&id, query, uuid)
	if err != nil && err == sql.ErrNoRows {
		return 0
	}
	return id
}

//serves users methods
func (users *Users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get users")
		res, err := users.GetUsersHandler(w, r)
		if err != nil && err.Error() == "not logged in" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodPut:
		fmt.Println("edit user")
		res, err := users.EditUserHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodDelete:
		fmt.Println("delete user")
		res, err := users.DeleteUserHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	default:
		return
	}
}

//serves user methods
func (user *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	switch {
	case r.Method == http.MethodGet:
		fmt.Println("get user")
		res, err := user.GetUserHandler(w, r)
		if err != nil && err.Error() == "not logged in" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil && err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	case r.Method == http.MethodPost:
		method := r.Header.Get("Method")
		var res []byte
		var err error
		if method == "login" {
			fmt.Println("post user login")
			res, err = user.LoginUserHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}

		if method == "logout" {
			fmt.Println("post user logout")
			res, err = user.LogoutUserHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
		}

		if method == "register" {
			fmt.Println("post user register")
			res, err = user.RegisterUserHandler(w, r)
			fmt.Println(res, err)
			if err != nil {
				if err.Error() == "display name already taken" {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(err.Error()))
					return
				}
				if err.Error() == "user already registered" {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(err.Error()))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
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
