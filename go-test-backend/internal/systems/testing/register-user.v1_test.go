package testing

import (
	"context"
	"test-backend/m/v2/internal/systems/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestRegisterUserV1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockUser := &user.RegisterUserV1Params{
		DisplayName: "test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@test.de",
		Birthday:    1234511,
		Password:    "secret",
	}

	mockResult := &user.RegisterUserV1Result{
		User: &user.RegisterUser{
			DisplayName: "test",
			FirstName:   "test",
			LastName:    "test",
			Email:       "test@test.de",
			Birthday:    1234511,
			Password:    "secret",
		},
	}

	context := context.Background()
	newUser := user.NewUserProvider(context, sqlxDB, nil, "sql")
	user := newUser.NewUser()

	mock.ExpectExec("^SELECT (.+)").WithArgs(12345).WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectExec("^INSERT (.+)").WithArgs(mockUser).WillReturnResult(sqlmock.NewResult(1, 1))
	// now we execute our method
	if err = user.RegisterUserV1(context, mockUser, mockResult); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
