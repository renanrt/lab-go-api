package model

import "strings"

const (
	AVALARA = "avalara"
	TAXJAR  = "taxjar"
)

// Represents a tax associated to a retailer.
//
// Saved in DB
//
// swagger:model tax
type Tax struct {
	// the id for this tax
	//
	// required: true
	// read only: true
	ID string `json:"id"`

	//retailer ID
	RetailerID string `json:"retailer_id"`

	// the name of the promotion
	//
	// required: true
	Name string `json:"name"`

	//tax ID saved in the monoliph
	VendTaxID string `json:"vend_tax_id"`

	//which provider where this tax was found (Avalara or TaxJar)
	Source SourceType `json:"source_id"`

	//the tax rate
	Rate float64 `json:"rate"`

	//this is the id of a parent tax. For example: a New York state tax
	ParentId string `json:"parent_id"`

	// tax type
	Type TaxType `json:"type"`
}

// SourceType Tax api providers
type SourceType string

//TaxType Type of identified tax
type TaxType string

const (
	// SourceTypeAvalara avalara
	SourceTypeAvalara SourceType = "avalara"
	// SourceTypeTaxjar taxjar
	SourceTypeTaxjar SourceType = "taxjar"
)

const (
	// TaxTypeState State tax
	TaxTypeState TaxType = "state"
	// TaxTypeCity City tax
	TaxTypeCity TaxType = "city"
	// TaxTypeCounty County tax
	TaxTypeCounty TaxType = "county"
	// TaxTypeSpecial Special type - usually it's a district type
	TaxTypeSpecial TaxType = "special"
)

//TaxGroup represents a group of taxes
type TaxGroup struct {
	TotalRate float64 `json:"total_rate"`
	//this field is used for confirming that the taxes were accepted
	RequestID string     `json:"request_id"`
	Rates     []*TaxRate `json:"rates"`
}

//TaxRate represents a single tax
type TaxRate struct {
	VendTaxID string  `json:"vend_tax_id"`
	Rate      float64 `json:"rate"`
	Name      string  `json:"name"`
	Type      TaxType `json:"type"`
}

//NewTaxRate Creates a new taxRate object
func NewTaxRate(t TaxType, name string, rate float64) *TaxRate {
	return &TaxRate{Name: name, Rate: rate, Type: t}
}

func (tr *TaxRate) ToTax() *Tax {
	return &Tax{Name: tr.Name, Rate: tr.Rate, Type: tr.Type, VendTaxID: tr.VendTaxID}
}

// GetTaxByType returns all tax rates that matches a type
func (tg *TaxGroup) GetTaxByType(t TaxType) []*TaxRate {
	r := []*TaxRate{}
	for _, tax := range tg.Rates {
		if IsSameType(tax.Type, t) {
			r = append(r, tax)
		}
	}
	return r
}

// ContainsTaxRate checks if the tax group contains a specific tax rate
func (tg *TaxGroup) ContainsTaxRate(tr *TaxRate) bool {
	for _, tax := range tg.Rates {
		if IsSameTaxRate(tax, tr) {
			return true
		}
	}
	return false
}

// AddTaxRate adds a taxRate to the tax group
func (tg *TaxGroup) AddTaxRate(taxRate *TaxRate) {
	tg.Rates = append(tg.Rates, taxRate)
}

// EllectParentId tries to determine what's the parentID for a taxGroup
// for US: it's the state tax
func (tg *TaxGroup) EllectParentId() string {
	if tg == nil {
		return ""
	}
	for _, tr := range tg.Rates {
		if tr.VendTaxID == "" {
			continue
		}
		if IsSameType(tr.Type, TaxTypeState) {
			return tr.VendTaxID
		}
	}
	return ""
}

// ContainsParentId checks if the parentID is one of the taxRates
func (tg *TaxGroup) ContainsParentId(parentID string) bool {
	if parentID == "" {
		return false
	}
	for _, t := range tg.Rates {
		if t.VendTaxID == parentID {
			return true
		}
	}
	return false
}

// FillTaxParentIds sets the vendTaxId on taxRate if the taxes are the same
// Example:
// The list 'taxesWithNoParent' contains all taxes from a retailer where the parentID is empty.
// Ideally, just State taxes have empty parentID's. Alll othe ones should have a parentID set
func (tg *TaxGroup) FillTaxParentIds(taxesWithNoParent []*Tax) {
	for _, apiTaxRate := range tg.GetTaxByType(TaxTypeState) {
		for _, dbTax := range taxesWithNoParent {
			if dbTax.ParentId != "" {
				continue
			}
			if IsSameTax(dbTax, apiTaxRate) {
				apiTaxRate.VendTaxID = dbTax.VendTaxID
			}
		}
	}
}

// IsSameTax determines if an apiTaxRate is the same from the one saved in DB
func IsSameTax(dbTax *Tax, apiTaxRate *TaxRate) bool {
	if !IsSameType(dbTax.Type, apiTaxRate.Type) {
		return false
	}
	return dbTax.Rate == apiTaxRate.Rate && strings.EqualFold(strings.TrimSpace(dbTax.Name), strings.TrimSpace(apiTaxRate.Name))
}

func IsSameTaxRate(taxRate1, taxRate2 *TaxRate) bool {
	if !IsSameType(taxRate1.Type, taxRate2.Type) {
		return false
	}
	return taxRate1.Rate == taxRate2.Rate && strings.EqualFold(strings.TrimSpace(taxRate1.Name), strings.TrimSpace(taxRate2.Name))
}

func IsAvalara(provider string) bool {
	return strings.ToLower(strings.TrimSpace(provider)) == AVALARA
}

func IsSameType(t1 TaxType, t2 TaxType) bool {

	s1 := strings.ToLower(strings.TrimSpace(string(t1)))
	s2 := strings.ToLower(strings.TrimSpace(string(t2)))
	return s1 == s2
}
