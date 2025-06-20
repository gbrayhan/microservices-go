package medicine

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
)

type IMedicineUseCase interface {
	GetData(page int64, limit int64, sortBy string, sortDirection string,
		filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error)
	GetByID(id int) (*domainMedicine.Medicine, error)
	Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error)
	GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error)
	GetAll() (*[]domainMedicine.Medicine, error)
}

type MedicineUseCase struct {
	MedicineRepository medicine.MedicineRepositoryInterface
}

func NewMedicineUseCase(medicineRepository medicine.MedicineRepositoryInterface) IMedicineUseCase {
	return &MedicineUseCase{
		MedicineRepository: medicineRepository,
	}
}

func (s *MedicineUseCase) GetData(page int64, limit int64, sortBy string, sortDirection string,
	filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	return s.MedicineRepository.GetData(page, limit, sortBy, sortDirection, filters, searchText, dateRangeFilters)
}

func (s *MedicineUseCase) GetByID(id int) (*domainMedicine.Medicine, error) {
	return s.MedicineRepository.GetByID(id)
}

func (s *MedicineUseCase) Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error) {
	return s.MedicineRepository.Create(medicine)
}

func (s *MedicineUseCase) GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	return s.MedicineRepository.GetByMap(medicineMap)
}

func (s *MedicineUseCase) Delete(id int) error {
	return s.MedicineRepository.Delete(id)
}

func (s *MedicineUseCase) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	return s.MedicineRepository.Update(id, medicineMap)
}

func (s *MedicineUseCase) GetAll() (*[]domainMedicine.Medicine, error) {
	return s.MedicineRepository.GetAll()
}
