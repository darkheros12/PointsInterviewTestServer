package internal

import (
	"PointsInterviewTestServer/internal/client"
	"PointsInterviewTestServer/internal/controllers"
	"PointsInterviewTestServer/internal/services"
	"PointsInterviewTestServer/internal/utils"
	"net/http"
)

// methodMux allows multiple HTTP methods for a single route.
func methodMux(handlers map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})
}

func NewRouter() *http.ServeMux {
	baseURL := "http://localhost:5001"

	// Initialize the HTTP client, cache, and tax calculator service
	cache := utils.NewMemoryCache()
	httpTaxClient := client.NewHTTPTaxClient(baseURL, cache)
	taxCalculatorService := services.NewTaxCalculator(httpTaxClient)

	mux := http.NewServeMux()
	mux.Handle("/tax-calculator/tax-year/", methodMux(map[string]http.Handler{
		http.MethodGet: controllers.GetTaxedIncomeWithBand(taxCalculatorService),
	}))

	return mux
}

// NewRouterWithBaseURL allows injecting a custom baseURL (for integration tests)
func NewRouterWithBaseURL(baseURL string) *http.ServeMux {
	cache := utils.NewMemoryCache()
	httpTaxClient := client.NewHTTPTaxClient(baseURL, cache)
	taxCalculatorService := services.NewTaxCalculator(httpTaxClient)

	mux := http.NewServeMux()
	mux.Handle("/tax-calculator/tax-year/", methodMux(map[string]http.Handler{
		http.MethodGet: controllers.GetTaxedIncomeWithBand(taxCalculatorService),
	}))

	return mux
}
