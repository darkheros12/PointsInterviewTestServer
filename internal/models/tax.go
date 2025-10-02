package models

type TaxBracket struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max,omitempty"`
	Rate float64 `json:"rate"`
}

type MockApiTaxResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

type TaxBandResult struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max,omitempty"`
	Rate    float64 `json:"rate"`
	Taxable float64 `json:"taxable"`
	Tax     float64 `json:"tax"`
}

type IncomeTaxResult struct {
	Year          int             `json:"year"`
	Income        float64         `json:"income"`
	TotalTax      float64         `json:"total_tax"`
	EffectiveRate float64         `json:"effective_rate"`
	Bands         []TaxBandResult `json:"bands"`
}
