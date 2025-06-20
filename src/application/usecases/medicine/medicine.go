package medicine

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"go.uber.org/zap"
)

type IMedicineUseCase interface {
	GetByID(id int) (*medicineDomain.Medicine, error)
	Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error)
	GetAll() (*[]medicineDomain.Medicine, error)
	SearchPaginated(filters domain.DataFilters) (*medicineDomain.SearchResultMedicine, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
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

func (s *MedicineUseCase) GetByID(id int) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Getting medicine by ID", zap.Int("id", id))
	return s.medicineRepository.GetByID(id)
}

func (s *MedicineUseCase) Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
	s.Logger.Info("Creating new medicine", zap.String("name", medicine.Name))
	return s.medicineRepository.Create(medicine)
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

func (s *MedicineUseCase) SearchPaginated(filters domain.DataFilters) (*medicineDomain.SearchResultMedicine, error) {
	s.Logger.Info("Searching medicines with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.medicineRepository.SearchPaginated(filters)
}

func (s *MedicineUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching medicines by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.medicineRepository.SearchByProperty(property, searchText)
}
