// Package medicine provides the use case for medicine
package medicine

import (
	"github.com/gbrayhan/microservices-go/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/domain/medicine"
	medicineRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/medicine"
)

// Service is a struct that contains the repository implementation for medicine use case
type Service struct {
	MedicineRepository medicineRepository.Repository
}

var _ medicineDomain.Service = &Service{}

// GetData is a function that returns all medicines
func (s *Service) GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error) {
	return s.MedicineRepository.GetAll(page, limit, sortBy, sortDirection, filters, searchText, dateRangeFilters)
}

// GetByID is a function that returns a medicine by id
func (s *Service) GetByID(id int) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.GetByID(id)
}

// Create is a function that creates a medicine
func (s *Service) Create(medicine *medicineDomain.NewMedicine) (*medicineDomain.Medicine, error) {
	medicineModel := medicine.ToDomainMapper()
	return s.MedicineRepository.Create(medicineModel)
}

// GetByMap is a function that returns a medicine by map
func (s *Service) GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.GetOneByMap(medicineMap)
}

// Delete is a function that deletes a medicine by id
func (s *Service) Delete(id int) error {
	return s.MedicineRepository.Delete(id)
}

// Update is a function that updates a medicine by id
func (s *Service) Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.Update(id, medicineMap)
}
