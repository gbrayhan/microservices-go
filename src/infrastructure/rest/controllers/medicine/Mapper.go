package medicine

import domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"

func domainToResponseMapper(m *domainMedicine.Medicine) *ResponseMedicine {
	return &ResponseMedicine{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		EanCode:     m.EanCode,
		Laboratory:  m.Laboratory,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(m *[]domainMedicine.Medicine) *[]ResponseMedicine {
	res := make([]ResponseMedicine, len(*m))
	for i, med := range *m {
		res[i] = *domainToResponseMapper(&med)
	}
	return &res
}
