package medicine

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
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
	require.NoError(t, err)
	return loggerInstance
}

func TestTableName(t *testing.T) {
	medicine := &Medicine{}
	assert.Equal(t, "medicines", medicine.TableName())
}

func TestNewMedicineRepository(t *testing.T) {
	db, _, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	assert.NotNil(t, repo)
}

func TestToDomainMapper(t *testing.T) {
	now := time.Now()
	medicine := &Medicine{
		ID:          1,
		Name:        "Test Medicine",
		Description: "Test Description",
		EANCode:     "1234567890123",
		Laboratory:  "Test Lab",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	domainMedicine := medicine.toDomainMapper()

	assert.Equal(t, medicine.ID, domainMedicine.ID)
	assert.Equal(t, medicine.Name, domainMedicine.Name)
	assert.Equal(t, medicine.Description, domainMedicine.Description)
	assert.Equal(t, medicine.EANCode, domainMedicine.EanCode)
	assert.Equal(t, medicine.Laboratory, domainMedicine.Laboratory)
	assert.Equal(t, medicine.CreatedAt, domainMedicine.CreatedAt)
	assert.Equal(t, medicine.UpdatedAt, domainMedicine.UpdatedAt)
}

func TestArrayToDomainMapper(t *testing.T) {
	now := time.Now()
	medicines := []Medicine{
		{
			ID:          1,
			Name:        "Medicine 1",
			Description: "Description 1",
			EANCode:     "1234567890123",
			Laboratory:  "Lab 1",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Name:        "Medicine 2",
			Description: "Description 2",
			EANCode:     "1234567890124",
			Laboratory:  "Lab 2",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	domainMedicines := arrayToDomainMapper(&medicines)

	assert.Len(t, *domainMedicines, 2)
	assert.Equal(t, medicines[0].ID, (*domainMedicines)[0].ID)
	assert.Equal(t, medicines[1].ID, (*domainMedicines)[1].ID)
}

// Integration-style tests using in-memory SQLite database
func TestRepository_GetAll(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "ean_code", "laboratory"}).
		AddRow(1, "Medicine 1", "Description 1", "1234567890123", "Lab 1").
		AddRow(2, "Medicine 2", "Description 2", "1234567890124", "Lab 2")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "medicines"`)).WillReturnRows(rows)
	medicines, err := repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, medicines)
	assert.Len(t, *medicines, 2)
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "ean_code", "laboratory"}).
		AddRow(1, "Medicine 1", "Description 1", "1234567890123", "Lab 1")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "medicines" WHERE id = $1 ORDER BY "medicines"."id" LIMIT $2`)).
		WithArgs(1, 1).WillReturnRows(rows)
	medicine, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, medicine)
	assert.Equal(t, "Medicine 1", medicine.Name)
}

func TestRepository_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	domainM := &medicineDomain.Medicine{
		Name:        "New Medicine",
		Description: "New Description",
		EanCode:     "1234567890125",
		Laboratory:  "New Lab",
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "medicines"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	medicine, err := repo.Create(domainM)
	assert.NoError(t, err)
	assert.NotNil(t, medicine)
	assert.Equal(t, "New Medicine", medicine.Name)
}

func TestRepository_Update(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "medicines" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "ean_code", "laboratory"}).
		AddRow(1, "Updated Medicine", "Updated Description", "1234567890123", "Updated Lab")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "medicines" WHERE id = $1 AND "medicines"."id" = $2 ORDER BY "medicines"."id" LIMIT $3`)).
		WithArgs(1, 1, 1).WillReturnRows(rows)
	medicine, err := repo.Update(1, map[string]any{"name": "Updated Medicine"})
	assert.NoError(t, err)
	assert.NotNil(t, medicine)
	assert.Equal(t, "Updated Medicine", medicine.Name)
}

func TestRepository_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewMedicineRepository(db, logger)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "medicines" WHERE "medicines"."id" = $1`)).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.Delete(1)
	assert.NoError(t, err)
}

func TestColumnsMedicineMapping(t *testing.T) {
	assert.Equal(t, "id", ColumnsMedicineMapping["id"])
	assert.Equal(t, "name", ColumnsMedicineMapping["name"])
	assert.Equal(t, "description", ColumnsMedicineMapping["description"])
	assert.Equal(t, "ean_code", ColumnsMedicineMapping["eanCode"])
	assert.Equal(t, "laboratory", ColumnsMedicineMapping["laboratory"])
	assert.Equal(t, "created_at", ColumnsMedicineMapping["createdAt"])
	assert.Equal(t, "updated_at", ColumnsMedicineMapping["updatedAt"])
}
