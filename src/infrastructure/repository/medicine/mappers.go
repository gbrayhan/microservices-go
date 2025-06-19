package medicine

import (
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
)

func (m *Medicine) toDomainMapper() *domainMedicine.Medicine {
	return &domainMedicine.Medicine{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		EanCode:     m.EANCode,
		Laboratory:  m.Laboratory,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func arrayToDomainMapper(medicines *[]Medicine) *[]domainMedicine.Medicine {
	medicinesDomain := make([]domainMedicine.Medicine, len(*medicines))
	for i, medicine := range *medicines {
		medicinesDomain[i] = *medicine.toDomainMapper()
	}
	return &medicinesDomain
}
