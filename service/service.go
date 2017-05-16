package service

import "github.com/renanrt/lab-go-api/model"

// TaxService represents the tax recommendation service
type TaxService struct {
}

func (service *TaxService) GetTaxesForAddress(provider, retailerId, country, state, city, zipcode string, street string) (*model.TaxGroup, error) {
	taxGroup := &model.TaxGroup{}
	taxGroup.TotalRate = 1.0
	return taxGroup, nil
}
