package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this will test that the method FillTaxParentIds sets the parentID properly
func TestSetParentId(t *testing.T) {
	tg := TaxGroup{}
	stateTax := NewTaxRate(TaxTypeState, "California", 0.0065)
	tg.AddTaxRate(stateTax)
	tg.AddTaxRate(NewTaxRate(TaxTypeCity, "Santa Monica", 0.001))

	dbStateTax := &Tax{Name: "California", Rate: 0.0065, Type: TaxTypeState, VendTaxID: "myVendTaxId"}
	dbCityTax := &Tax{Name: "NYC", Rate: 0.0065, Type: TaxTypeCity, VendTaxID: "myVendCityTaxId", ParentId: dbStateTax.VendTaxID}

	tg.FillTaxParentIds([]*Tax{dbStateTax, dbCityTax})
	assert.NotEmpty(t, stateTax.VendTaxID)
	assert.True(t, tg.ContainsParentId(stateTax.VendTaxID))
	assert.False(t, tg.ContainsParentId(dbCityTax.VendTaxID))
}

// this will test that the method FillTaxParentIds sets the parentID properly
func TestContainsTaxRate(t *testing.T) {
	tg := TaxGroup{}
	stateTax := NewTaxRate(TaxTypeState, "California", 0.0065)
	cityTax := NewTaxRate(TaxTypeCity, "Santa Monica", 0.001)
	ghostCityTax := NewTaxRate(TaxTypeCity, "AnotherCity", 0.001)
	tg.AddTaxRate(stateTax)
	tg.AddTaxRate(cityTax)

	assert.True(t, tg.ContainsTaxRate(stateTax))
	assert.True(t, tg.ContainsTaxRate(cityTax))
	assert.False(t, tg.ContainsTaxRate(ghostCityTax))
}

func TestEllectParentID(t *testing.T) {
	tg := TaxGroup{}
	stateTax := NewTaxRate(TaxTypeState, "California", 0.0065)
	stateTax.VendTaxID = "vend-tax-id"
	tg.AddTaxRate(stateTax)
	tg.AddTaxRate(NewTaxRate(TaxTypeCity, "Santa Monica", 0.001))

	assert.NotEmpty(t, tg.EllectParentId())
	assert.Equal(t, stateTax.VendTaxID, tg.EllectParentId())
}

// tests when it's the same tax
func TestIsSameTax(t *testing.T) {
	stateTax := NewTaxRate(TaxTypeState, "California", 0.0065)
	dbTax := &Tax{Name: "California", Rate: 0.0065, Type: TaxTypeState, VendTaxID: "myVendTaxId"}
	assert.True(t, IsSameTax(dbTax, stateTax))

	dbTax.Type = "State"
	assert.True(t, IsSameTax(dbTax, stateTax))

	dbTax.Name = " caLiFornia "
	assert.True(t, IsSameTax(dbTax, stateTax))

	dbTax.Name = " california _"
	assert.False(t, IsSameTax(dbTax, stateTax))
}

// tests when two taxTypes are the same
func TestIsSameTaxType(t *testing.T) {
	var t1, t2 TaxType
	t1 = TaxTypeState
	t2 = TaxTypeState
	assert.True(t, IsSameType(t1, t2))

	t2 = " State "
	assert.True(t, IsSameType(t1, t2))
}
