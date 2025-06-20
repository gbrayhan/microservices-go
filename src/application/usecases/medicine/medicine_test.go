package medicine

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
)

type mockMedicineService struct {
	getDataFn func(page int64, limit int64, sortBy string, sortDirection string,
		filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error)
	getByIDFn  func(id int) (*medicineDomain.Medicine, error)
	createFn   func(m *medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	getByMapFn func(m map[string]any) (*medicineDomain.Medicine, error)
	deleteFn   func(id int) error
	updateFn   func(id int, m map[string]any) (*medicineDomain.Medicine, error)
	getAllFn   func() (*[]medicineDomain.Medicine, error)
}

func (m *mockMedicineService) GetData(page int64, limit int64, sortBy string, sortDirection string,
	filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error) {
	return m.getDataFn(page, limit, sortBy, sortDirection, filters, searchText, dateRangeFilters)
}

func (m *mockMedicineService) GetByID(id int) (*medicineDomain.Medicine, error) {
	return m.getByIDFn(id)
}

func (m *mockMedicineService) Create(med *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
	return m.createFn(med)
}

func (m *mockMedicineService) GetByMap(med map[string]any) (*medicineDomain.Medicine, error) {
	return m.getByMapFn(med)
}

func (m *mockMedicineService) Delete(id int) error {
	return m.deleteFn(id)
}

func (m *mockMedicineService) Update(id int, med map[string]any) (*medicineDomain.Medicine, error) {
	return m.updateFn(id, med)
}

func (m *mockMedicineService) GetAll() (*[]medicineDomain.Medicine, error) {
	return m.getAllFn()
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestMedicineUseCase(t *testing.T) {
	mockRepo := &mockMedicineService{}
	loggerInstance := setupLogger(t)
	useCase := NewMedicineUseCase(mockRepo, loggerInstance)

	mockRepo.getDataFn = func(page int64, limit int64, sortBy string, sortDirection string,
		filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error) {
		if page == 0 {
			return nil, errors.New("error in repository")
		}
		return &medicineDomain.DataMedicine{
			Data:  &[]medicineDomain.Medicine{{ID: 1, Name: "Med1"}},
			Total: 1,
		}, nil
	}
	_, err := useCase.GetData(0, 10, "", "", nil, "", nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
	medData, err := useCase.GetData(1, 10, "", "", nil, "", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if medData == nil || len(*medData.Data) != 1 {
		t.Error("expected one medicine in data")
	}

	mockRepo.getByIDFn = func(id int) (*medicineDomain.Medicine, error) {
		if id == 123 {
			return &medicineDomain.Medicine{ID: 123}, nil
		}
		return nil, errors.New("not found")
	}
	_, err = useCase.GetByID(999)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}
	med, err := useCase.GetByID(123)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if med.ID != 123 {
		t.Error("expected medicine ID=123")
	}

	mockRepo.createFn = func(m *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
		if m.Name == "" {
			return nil, errors.New("validation error")
		}
		m.ID = 999
		return m, nil
	}
	_, err = useCase.Create(&medicineDomain.Medicine{Name: ""})
	if err == nil {
		t.Error("expected create error on empty name")
	}
	newMed, err := useCase.Create(&medicineDomain.Medicine{Name: "Aspirin"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if newMed.ID != 999 {
		t.Error("expected created medicine ID=999")
	}

	mockRepo.getByMapFn = func(mm map[string]any) (*medicineDomain.Medicine, error) {
		if mm["ean"] == "not_found" {
			return nil, errors.New("not found")
		}
		return &medicineDomain.Medicine{ID: 55}, nil
	}
	_, err = useCase.GetByMap(map[string]any{"ean": "not_found"})
	if err == nil {
		t.Error("expected error, got nil")
	}
	mg, err := useCase.GetByMap(map[string]any{"ean": "valid"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if mg.ID != 55 {
		t.Error("expected ID=55")
	}

	mockRepo.deleteFn = func(id int) error {
		if id == 1010 {
			return nil
		}
		return errors.New("cannot delete")
	}
	err = useCase.Delete(100)
	if err == nil {
		t.Error("expected error, got nil")
	}
	err = useCase.Delete(1010)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mockRepo.updateFn = func(id int, mm map[string]any) (*medicineDomain.Medicine, error) {
		if id != 1000 {
			return nil, errors.New("not found for update")
		}
		return &medicineDomain.Medicine{ID: 1000, Name: "UpdatedName"}, nil
	}
	_, err = useCase.Update(999, map[string]any{"name": "whatever"})
	if err == nil {
		t.Error("expected error, got nil")
	}
	updated, err := useCase.Update(1000, map[string]any{"name": "NewName"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updated.Name != "UpdatedName" {
		t.Error("expected updated name to be UpdatedName")
	}

	mockRepo.getAllFn = func() (*[]medicineDomain.Medicine, error) {
		return &[]medicineDomain.Medicine{
			{ID: 1, Name: "M1"}, {ID: 2, Name: "M2"},
		}, nil
	}
	meds, err := useCase.GetAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if meds == nil || len(*meds) != 2 {
		t.Error("expected 2 medicines from GetAll()")
	}
}

func TestNewMedicineUseCase(t *testing.T) {
	mockRepo := &mockMedicineService{}
	loggerInstance := setupLogger(t)
	uc := NewMedicineUseCase(mockRepo, loggerInstance)
	if reflect.TypeOf(uc).String() != "*medicine.MedicineUseCase" {
		t.Error("expected *medicine.MedicineUseCase type")
	}
}
