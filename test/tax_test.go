package test

import (
	"PointsInterviewTestServer/internal"
	"PointsInterviewTestServer/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// startMockTaxAPI starts a mock HTTP server that returns the given status and response.
func startMockTaxAPI(t *testing.T, status int, resp models.MockApiTaxResponse) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestIntegration_ValidRequest(t *testing.T) {
	mockResp := models.MockApiTaxResponse{
		TaxBrackets: []models.TaxBracket{
			{Min: 0, Max: 50000, Rate: 0.15},
		},
	}
	mockAPI := startMockTaxAPI(t, 200, mockResp)

	router := internal.NewRouterWithBaseURL(mockAPI.URL)
	req := httptest.NewRequest(http.MethodGet, "/tax-calculator/tax-year/?year=2022&salary=75000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got models.IncomeTaxResult
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.TotalTax == 0 {
		t.Errorf("expected nonzero tax, got 0")
	}
}

func TestIntegration_MissingParams(t *testing.T) {
	mockResp := models.MockApiTaxResponse{
		TaxBrackets: []models.TaxBracket{
			{Min: 0, Max: 50000, Rate: 0.15},
		},
	}
	mockAPI := startMockTaxAPI(t, 400, mockResp)

	router := internal.NewRouterWithBaseURL(mockAPI.URL)
	req := httptest.NewRequest(http.MethodGet, "/tax-calculator/tax-year/2022", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestIntegration_UpstreamError(t *testing.T) {
	mockAPI := startMockTaxAPI(t, 500, models.MockApiTaxResponse{})

	router := internal.NewRouterWithBaseURL(mockAPI.URL)
	req := httptest.NewRequest(http.MethodGet, "/tax-calculator/tax-year/?year=2022&salary=75000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", w.Code)
	}
}
