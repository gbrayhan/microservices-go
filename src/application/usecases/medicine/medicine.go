package medicine

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"go.uber.org/zap"
)

type IMedicineUseCase interface {
	GetData(page int64, limit int64, sortBy string, sortDirection string,
		filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error)
	GetByID(id int) (*medicineDomain.Medicine, error)
	Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error)
	GetAll() (*[]medicineDomain.Medicine, error)
}

type MedicineUseCase struct {
	medicineRepository medicine.MedicineRepositoryInterface
	Logger             *logger.Logger
}

func NewMedicineUseCase(medicineRepository medicine.MedicineRepositoryInterface, loggerInstance *logger.Logger) IMedicineUseCase {
	return &MedicineUseCase{
		medicineRepository: medicineRepository,
		Logger:             loggerInstance,
	}
}

func (s *MedicineUseCase) GetData(page int64, limit int64, sortBy string, sortDirection string,
	filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error) {
	s.Logger.Info("Getting medicine data", zap.Int64("page", page), zap.Int64("limit", limit))
	return s.medicineRepository.GetData(page, limit, sortBy, sortDirection, filters, searchText, dateRangeFilters)
}

func (s *MedicineUseCase) GetByID(id int) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Getting medicine by ID", zap.Int("id", id))
	return s.medicineRepository.GetByID(id)
}

func (s *MedicineUseCase) Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Creating new medicine", zap.String("name", medicine.Name))
	return s.medicineRepository.Create(medicine)
}

func (s *MedicineUseCase) GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Getting medicine by map", zap.Any("medicineMap", medicineMap))
	return s.medicineRepository.GetByMap(medicineMap)
}

func (s *MedicineUseCase) Delete(id int) error {
	s.Logger.Info("Deleting medicine", zap.Int("id", id))
	return s.medicineRepository.Delete(id)
}

func (s *MedicineUseCase) Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Updating medicine", zap.Int("id", id))
	return s.medicineRepository.Update(id, medicineMap)
}

func (s *MedicineUseCase) GetAll() (*[]medicineDomain.Medicine, error) {
	s.Logger.Info("Getting all medicines")
	return s.medicineRepository.GetAll()
}
