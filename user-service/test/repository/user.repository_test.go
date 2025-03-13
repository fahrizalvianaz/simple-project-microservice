package repository_test

import (
	"bookstore-framework/internal/users"
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestUserRepository_Success(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gormDB, mock := setupMockDB(t)
	repo := users.NewUserRepository(gormDB)
	t.Run("Register", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		testUser := &users.User{
			Username: "testuser",
			Email:    "test@gmail.com",
			Password: "password123",
		}

		result, err := repo.Register(context.Background(), testUser)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "testuser", result.Username)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("FindUserByUsername", func(t *testing.T) {
		username := "test"
		columns := []string{"id", "username", "name", "email", "password", "created_at", "modified_at", "deleted_at"}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL`)).
			WithArgs(username, 1).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(1, "test", "testuser", "test@example.com", "hashedpassword", time.Now(), time.Now(), nil))

		user, err := repo.FindUserByUsername(context.Background(), username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("FindUserByID", func(t *testing.T) {
		idUser := 1
		columns := []string{"id", "username", "name", "email", "password", "created_at", "modified_at", "deleted_at"}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).
			WithArgs(idUser, 1).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(1, "test", "testuser", "test@example.com", "hashedpassword", time.Now(), time.Now(), nil))
		user, err := repo.FindUserByID(context.Background(), uint(idUser))
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int(idUser), int(user.ID))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestUserRepository_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gormDB, mock := setupMockDB(t)
	repo := users.NewUserRepository(gormDB)

	t.Run("Register", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnError(errors.New("Error database"))
		mock.ExpectRollback()

		testUser := &users.User{
			Username: "testuser",
			Email:    "test@gmail.com",
			Password: "password123",
		}

		result, err := repo.Register(context.Background(), testUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Error database", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("FindUserByUsername", func(t *testing.T) {
		username := "test"
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL`)).
			WithArgs(username, 1).
			WillReturnError(errors.New("User not found"))

		user, err := repo.FindUserByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "User not found", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("FindUserByID", func(t *testing.T) {
		idUser := 1

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).
			WithArgs(idUser, 1).
			WillReturnError(errors.New("User not found"))

		user, err := repo.FindUserByID(context.Background(), uint(idUser))
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "User not found", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})
}
