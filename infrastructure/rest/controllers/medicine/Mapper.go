package medicine

import medicineDomain "github.com/gbrayhan/microservices-go/domain/medicine"

func domainToResponseMapper(clientDomain *medicineDomain.Medicine) (createClientResponse *ResponseMedicine) {
	createClientResponse = &ResponseMedicine{
		ID:          clientDomain.ID,
		Name:        clientDomain.Name,
		Description: clientDomain.Description,
		EanCode:     clientDomain.EanCode,
		Laboratory:  clientDomain.Laboratory,
		CreatedAt:   clientDomain.CreatedAt,
		UpdatedAt:   clientDomain.UpdatedAt}

	return
}

func arrayDomainToResponseMapper(clientsDomain *[]medicineDomain.Medicine) *[]ResponseMedicine {
	clientsResponse := make([]ResponseMedicine, len(*clientsDomain))
	for i, client := range *clientsDomain {
		clientsResponse[i] = *domainToResponseMapper(&client)
	}
	return &clientsResponse
}
