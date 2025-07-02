package user

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)
	cleanup := func() { db.Close() }
	return gormDB, mock, cleanup
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestTableName(t *testing.T) {
	u := &User{}
	assert.Equal(t, "users", u.TableName())
}

func TestNewUserRepository(t *testing.T) {
	db, _, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	assert.NotNil(t, repo)
}

func TestToDomainMapper(t *testing.T) {
	u := &User{
		ID:        1,
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	d := u.toDomainMapper()
	assert.Equal(t, u.UserName, d.UserName)
}

func TestFromDomainMapper(t *testing.T) {
	d := &domainUser.User{
		ID:        1,
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	u := fromDomainMapper(d)
	assert.Equal(t, d.UserName, u.UserName)
}

func TestArrayToDomainMapper(t *testing.T) {
	arr := &[]User{{ID: 1, UserName: "A"}, {ID: 2, UserName: "B"}}
	d := arrayToDomainMapper(arr)
	assert.Len(t, *d, 2)
	assert.Equal(t, "A", (*d)[0].UserName)
}

func TestRepository_GetAll(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password"}).
		AddRow(1, "user1", "a@a.com", "A", "B", true, "hash1").
		AddRow(2, "user2", "b@b.com", "C", "D", false, "hash2")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(rows)
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, *users, 2)
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password"}).
		AddRow(1, "user1", "a@a.com", "A", "B", true, "hash1")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).WillReturnRows(rows)
	user, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.UserName)
	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(2, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password"}))
	user, err = repo.GetByID(2)
	assert.Error(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 0, user.ID) // Should be zero value
}

func TestRepository_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	domainU := &domainUser.User{
		UserName:     "user1",
		Email:        "a@a.com",
		FirstName:    "A",
		LastName:     "B",
		Status:       true,
		HashPassword: "hash1",
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	user, err := repo.Create(domainU)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.UserName)
}

func TestRepository_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.Delete(1)
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	err = repo.Delete(2)
	assert.Error(t, err)
}

func TestRepository_GetByEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)

	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password"}).
		AddRow(1, "user1", email, "A", "B", true, "hash1")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).WillReturnRows(rows)
	user, err := repo.GetByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	// Not found
	emailNotFound := "notfound@example.com"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(emailNotFound, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password"}))
	user, err = repo.GetByEmail(emailNotFound)
	assert.Error(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 0, user.ID) // Should be zero value
}

// The following tests need refactoring to use sqlmock or should be moved to integration:
// TestRepository_GetOneByMap
// TestRepository_Update
// TestRepository_Create_DuplicateEmail
// TestRepository_ErrorCases
// TestRepository_GetOneByMap_WithFilters
// TestRepository_Update_WithMultipleFields
//
// If you want me to refactor these as well, let me know and I'll do them one by one.
