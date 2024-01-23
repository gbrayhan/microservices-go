// Package medicine provides the use case for medicine
package medicine

import (
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
)

func (n *NewMedicine) toDomainMapper() *domainMedicine.Medicine {
	return &domainMedicine.Medicine{
		Name:        n.Name,
		Description: n.Description,
		EanCode:     n.EANCode,
		Laboratory:  n.Laboratory,
	}
}
