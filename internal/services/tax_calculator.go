package services

import (
	"PointsInterviewTestServer/internal/client"
	"PointsInterviewTestServer/internal/models"
	"context"
	"math"
)

type ITaxCalculator interface {
	CalculateTax(ctx context.Context, year int, income float64) (*models.IncomeTaxResult, error)
	GetYearTaxBand(ctx context.Context, year int, income float64) (models.TaxBandResult, error)
}

type TaxCalculator struct {
	provider client.TaxRateProvider
}

func NewTaxCalculator(provider client.TaxRateProvider) *TaxCalculator {
	return &TaxCalculator{provider: provider}
}

// GetYearTaxBand To only get the year tax band. Not needed now but for future purpose.
func (taxCalc *TaxCalculator) GetYearTaxBand(ctx context.Context, year int, income float64) (models.TaxBandResult, error) {
	brackets, err := taxCalc.provider.GetTaxBrackets(ctx, year)
	if err != nil {
		return models.TaxBandResult{}, err
	}

	var bands []models.TaxBandResult
	var model models.TaxBandResult

	// Bands tax calculation for each bracket
	for _, value := range brackets {
		if income <= value.Min {
			model = models.TaxBandResult{
				Min:  value.Min,
				Max:  value.Max,
				Rate: value.Rate,
			}
			bands = append(bands, model)
			continue
		}
	}

	return model, nil
}

func (taxCalc *TaxCalculator) CalculateTax(ctx context.Context, year int, income float64) (*models.IncomeTaxResult, error) {
	brackets, err := taxCalc.provider.GetTaxBrackets(ctx, year)
	if err != nil {
		return nil, err
	}

	totalTax := 0.0
	var bands []models.TaxBandResult
	var model models.TaxBandResult

	// Bands tax calculation for each bracket
	for _, value := range brackets {
		if income <= value.Min {
			model = models.TaxBandResult{
				Min:  value.Min,
				Max:  value.Max,
				Rate: value.Rate,
			}
			bands = append(bands, model)
			continue
		}

		var taxable float64
		if value.Max > 0 {
			taxable = math.Min(income, value.Max) - value.Min
		} else {
			taxable = income - value.Min
		}
		if taxable < 0 {
			taxable = 0
		}
		tax := taxable * value.Rate
		totalTax += tax
		model = models.TaxBandResult{
			Min:     value.Min,
			Max:     value.Max,
			Rate:    value.Rate,
			Taxable: roundToTwo(taxable),
			Tax:     roundToTwo(tax),
		}
		bands = append(bands, model)
	}

	effectiveRate := 0.0
	if income > 0 {
		effectiveRate = totalTax / income
	}

	res := &models.IncomeTaxResult{
		Year:          year,
		Income:        roundToTwo(income),
		TotalTax:      roundToTwo(totalTax),
		EffectiveRate: roundToTwo(effectiveRate),
		Bands:         bands,
	}

	return res, nil
}

func roundToTwo(input float64) float64 {
	return math.Round(input*100) / 100
}
