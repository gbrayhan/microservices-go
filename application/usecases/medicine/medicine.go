// Package medicine provides the use case for medicine
package medicine

import (
	medicineDomain "github.com/gbrayhan/microservices-go/domain/medicine"
	medicineRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/medicine"
)

// Service is a struct that contains the repository implementation for medicine use case
type Service struct {
	MedicineRepository medicineRepository.Repository
}

// GetAll is a function that returns all medicines
func (s *Service) GetAll(page int64, limit int64) (*PaginationResultMedicine, error) {
	all, err := s.MedicineRepository.GetAll(page, limit)
	if err != nil {
		return nil, err
	}

	return &PaginationResultMedicine{
		Data:       all.Data,
		Total:      all.Total,
		Limit:      all.Limit,
		Current:    all.Current,
		NextCursor: all.NextCursor,
		PrevCursor: all.PrevCursor,
		NumPages:   all.NumPages,
	}, nil
}

// GetByID is a function that returns a medicine by id
func (s *Service) GetByID(id int) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.GetByID(id)
}

// Create is a function that creates a medicine
func (s *Service) Create(medicine *NewMedicine) (*medicineDomain.Medicine, error) {
	medicineModel := medicine.toDomainMapper()
	return s.MedicineRepository.Create(medicineModel)
}

// GetByMap is a function that returns a medicine by map
func (s *Service) GetByMap(medicineMap map[string]interface{}) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.GetOneByMap(medicineMap)
}

// Delete is a function that deletes a medicine by id
func (s *Service) Delete(id int) error {
	return s.MedicineRepository.Delete(id)
}

// Update is a function that updates a medicine by id
func (s *Service) Update(id int, medicineMap map[string]interface{}) (*medicineDomain.Medicine, error) {
	return s.MedicineRepository.Update(id, medicineMap)
}
