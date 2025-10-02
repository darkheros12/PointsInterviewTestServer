package services

import (
	"PointsInterviewTestServer/internal/models"
	"context"
	"reflect"
	"testing"
)

type fakeTaxClient struct {
	brackets []models.TaxBracket
	err      error
}

func (f *fakeTaxClient) GetTaxBrackets(ctx context.Context, year int) ([]models.TaxBracket, error) {
	return f.brackets, f.err
}

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name      string
		income    float64
		brackets  []models.TaxBracket
		wantTotal float64
		wantBands []models.TaxBandResult
		wantEff   float64
		wantErr   bool
	}{
		{
			"one bracket",
			40000,
			[]models.TaxBracket{{Min: 0, Max: 50000, Rate: 0.15}},
			6000,
			[]models.TaxBandResult{{Min: 0, Max: 50000, Rate: 0.15, Taxable: 40000, Tax: 6000}},
			0.15,
			false,
		},
		{
			"multiple brackets",
			60000,
			[]models.TaxBracket{{Min: 0, Max: 50000, Rate: 0.10}, {Min: 50000, Max: 100000, Rate: 0.20}},
			7000,
			[]models.TaxBandResult{
				{Min: 0, Max: 50000, Rate: 0.1, Taxable: 50000, Tax: 5000},
				{Min: 50000, Max: 100000, Rate: 0.2, Taxable: 10000, Tax: 2000}},
			roundToTwo(7000.00 / 60000.00),
			false,
		},
		{
			"open ended bracket",
			120000,
			[]models.TaxBracket{{Min: 0, Max: 50000, Rate: 0.10}, {Min: 50000, Max: 100000, Rate: 0.20}, {Min: 100000, Max: 0, Rate: 0.30}},
			21000,
			[]models.TaxBandResult{
				{Min: 0, Max: 50000, Rate: 0.1, Taxable: 50000, Tax: 5000},
				{Min: 50000, Max: 100000, Rate: 0.2, Taxable: 50000, Tax: 10000},
				{Min: 100000, Max: 0, Rate: 0.3, Taxable: 20000, Tax: 6000}},
			roundToTwo(21000.00 / 120000.00),
			false,
		},
		{
			"error from client",
			50000,
			nil,
			0,
			nil,
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewTaxCalculator(&fakeTaxClient{brackets: tt.brackets, err: func() error {
				if tt.wantErr && tt.brackets == nil {
					return context.Canceled
				}
				return nil
			}()})
			res, err := calc.CalculateTax(context.Background(), 2022, tt.income)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.TotalTax != tt.wantTotal {
				t.Errorf("got total tax %.2f, want %.2f", res.TotalTax, tt.wantTotal)
			}
			if !reflect.DeepEqual(res.Bands, tt.wantBands) {
				t.Errorf("got bands %+v, want %+v", res.Bands, tt.wantBands)
			}
			if res.EffectiveRate != tt.wantEff {
				t.Errorf("got effective rate %.4f, want %.4f", res.EffectiveRate, tt.wantEff)
			}
		})
	}
}
