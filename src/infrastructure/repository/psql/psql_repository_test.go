package psql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAuthService struct{}

func (m *mockAuthService) HashPassword(password string) (string, error) {
	return "hashed_" + password, nil
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "secretpassword",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expectedDSN := "host=localhost port=5432 user=postgres password=secretpassword dbname=testdb sslmode=disable TimeZone=America/Mexico_City"
	assert.Equal(t, expectedDSN, cfg.GetDSN())
}

func TestLoadDatabaseConfig_Success(t *testing.T) {
	// Set all required environment variables
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "dbuser")
	t.Setenv("DB_PASSWORD", "dbpass")
	t.Setenv("DB_NAME", "microservices")
	t.Setenv("DB_SSLMODE", "require")

	cfg, err := loadDatabaseConfig()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, "5433", cfg.Port)
	assert.Equal(t, "dbuser", cfg.User)
	assert.Equal(t, "dbpass", cfg.Password)
	assert.Equal(t, "microservices", cfg.DBName)
	assert.Equal(t, "require", cfg.SSLMode)
}

func TestLoadDatabaseConfig_MissingVariables(t *testing.T) {
	// Clear all DB environment variables
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_SSLMODE", "")

	cfg, err := loadDatabaseConfig()
	assert.Error(t, err)
	assert.Empty(t, cfg.Host)
	assert.Contains(t, err.Error(), "missing required database environment variables")
	assert.Contains(t, err.Error(), "DB_HOST")
	assert.Contains(t, err.Error(), "DB_PORT")
	assert.Contains(t, err.Error(), "DB_USER")
	assert.Contains(t, err.Error(), "DB_PASSWORD")
	assert.Contains(t, err.Error(), "DB_NAME")
	assert.Contains(t, err.Error(), "DB_SSLMODE")
}

func TestLoadDatabaseConfig_PartialMissing(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "")

	cfg, err := loadDatabaseConfig()
	assert.Error(t, err)
	assert.Empty(t, cfg.Host)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
	assert.Contains(t, err.Error(), "DB_SSLMODE")
	assert.NotContains(t, err.Error(), "DB_HOST")
}

func TestPSQLRepository_Setters(t *testing.T) {
	repo := NewRepository(nil, nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.DB)
	assert.Nil(t, repo.Logger)
	assert.Nil(t, repo.Auth)

	mockAuth := &mockAuthService{}
	repo.SetAuthService(mockAuth)
	assert.Equal(t, mockAuth, repo.Auth)

	hashed, err := repo.Auth.HashPassword("test1234")
	assert.NoError(t, err)
	assert.Equal(t, "hashed_test1234", hashed)

	repo.SetLogger(nil)
	assert.Nil(t, repo.Logger)
}

func TestPSQLRepository_LoadDBConfig(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "disable")

	repo := NewRepository(nil, nil)
	cfg, err := repo.LoadDBConfig()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "postgres", cfg.User)
	assert.Equal(t, "secret", cfg.Password)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "disable", cfg.SSLMode)
}
