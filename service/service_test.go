package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vend/princejohn/model"
)

func TestMergeTaxesAllNew(t *testing.T) {

	dbTaxes := []*model.Tax{}

	taxGroup := &model.TaxGroup{}
	// state tax
	taxGroup.AddTaxRate(californiaStateTaxRate())
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeCity, "SANTA MONICA", 0.005))
	taxGroup.AddTaxRate(losAngelesCountyTaxRate())
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeSpecial, "LOS ANGELES COUNTY DISTRICT TAX SP", 0.0015))
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeSpecial, "LOS ANGELES CO LOCAL TAX SL", 0.01))

	mergeTaxes(dbTaxes, taxGroup)
	assert.NotNil(t, taxGroup)
	assert.Equal(t, 5, len(taxGroup.Rates))

	for _, r := range taxGroup.Rates {
		assert.Empty(t, r.VendTaxID)
	}
}

func TestMergeTaxesReuseCountyAndState(t *testing.T) {
	vendStateTaxId := "any_vend_state_tax_id"
	vendCityTaxId := "any_vend_city_tax_id"
	dbTaxes := []*model.Tax{}
	dbTaxes = append(dbTaxes, &model.Tax{Name: "CALIFORNIA", Rate: 0.0625, ID: "0909", VendTaxID: vendStateTaxId, Type: model.TaxTypeState})
	dbTaxes = append(dbTaxes, &model.Tax{Name: "SANTA MONICA", Rate: 0.005, ID: "091", VendTaxID: vendCityTaxId, Type: model.TaxTypeCity, ParentId: vendStateTaxId})

	taxGroup := &model.TaxGroup{}
	// state tax
	taxGroup.AddTaxRate(californiaStateTaxRate())
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeCity, "SANTA MONICA", 0.005))
	taxGroup.AddTaxRate(losAngelesCountyTaxRate())
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeSpecial, "LOS ANGELES COUNTY DISTRICT TAX SP", 0.0015))
	taxGroup.AddTaxRate(model.NewTaxRate(model.TaxTypeSpecial, "LOS ANGELES CO LOCAL TAX SL", 0.01))

	// that's the testing method
	mergeTaxes(dbTaxes, taxGroup)
	assert.NotNil(t, taxGroup)
	assert.Equal(t, 5, len(taxGroup.Rates))

	// checking that state tax will be reused
	stateTaxes := taxGroup.GetTaxByType(model.TaxTypeState)
	assert.Equal(t, 1, len(stateTaxes))
	stateTax := stateTaxes[0]
	assert.NotEmpty(t, stateTax.VendTaxID)
	assert.Equal(t, vendStateTaxId, stateTax.VendTaxID)

	// checking that city tax will be reused
	cityTaxes := taxGroup.GetTaxByType(model.TaxTypeCity)
	assert.Equal(t, 1, len(cityTaxes))
	cityTax := cityTaxes[0]
	assert.NotEmpty(t, cityTax.VendTaxID)
	assert.Equal(t, vendCityTaxId, cityTax.VendTaxID)
}

func californiaStateTaxRate() *model.TaxRate {
	return model.NewTaxRate(model.TaxTypeState, "CALIFORNIA", 0.0625)
}

func losAngelesCountyTaxRate() *model.TaxRate {
	return model.NewTaxRate(model.TaxTypeCounty, "LOS ANGELES", 0.0025)
}
