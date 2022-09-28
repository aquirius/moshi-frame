package testing

import (
	"context"
	"test-backend/m/v2/internal/systems/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestDeleteUserV1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockUser := &user.DeleteUserV1Params{
		ID: "12345",
	}

	context := context.Background()
	newUsers := user.NewUsersProvider(context, sqlxDB, nil, "sql")
	users := newUsers.NewUsers()

	mock.ExpectExec("^DELETE (.+)").WithArgs(mockUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	// now we execute our method
	if err = users.DeleteUserV1(mockUser); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
