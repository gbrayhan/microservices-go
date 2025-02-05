package adapter

import (
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	medicineRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/medicine"
	medicineController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	"gorm.io/gorm"
)

func MedicineAdapter(db *gorm.DB) medicineController.IMedicineController {
	repository := medicineRepository.NewMedicineRepository(db)
	service := medicineUseCase.NewMedicineUseCase(repository)
	controller := medicineController.NewMedicineController(service)
	return controller
}
