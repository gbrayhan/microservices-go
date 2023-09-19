// Package medicine contains the repository implementation for the medicine entity
package medicine

import domainMedicine "github.com/gbrayhan/microservices-go/domain/medicine"

func (medicine *Medicine) toDomainMapper() *domainMedicine.Medicine {
	return &domainMedicine.Medicine{
		ID:          medicine.ID,
		Name:        medicine.Name,
		Description: medicine.Description,
		EanCode:     medicine.EANCode,
		Laboratory:  medicine.Laboratory,
		CreatedAt:   medicine.CreatedAt,
	}
}

func fromDomainMapper(medicine *domainMedicine.Medicine) *Medicine {
	return &Medicine{
		ID:          medicine.ID,
		Name:        medicine.Name,
		Description: medicine.Description,
		EANCode:     medicine.EanCode,
		Laboratory:  medicine.Laboratory,
		CreatedAt:   medicine.CreatedAt,
	}
}

func arrayToDomainMapper(medicines *[]Medicine) *[]domainMedicine.Medicine {
	medicinesDomain := make([]domainMedicine.Medicine, len(*medicines))
	for i, medicine := range *medicines {
		medicinesDomain[i] = *medicine.toDomainMapper()
	}

	return &medicinesDomain
}
