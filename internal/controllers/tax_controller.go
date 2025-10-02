package controllers

import (
	"PointsInterviewTestServer/internal/services"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

func GetTaxedIncomeWithBand(taxCalculator services.ITaxCalculator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		year, income, err := requestValidation(w, r)
		if err != nil {
			return
		}

		res, err := taxCalculator.CalculateTax(r.Context(), year, income)
		if err != nil {
			log.Printf("cannot calulate the tax income error: %v", err)
			http.Error(w, fmt.Sprintf("upstream error: %v", err), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})
}

func requestValidation(w http.ResponseWriter, r *http.Request) (year int, income float64, err error) {
	salaryStr := r.URL.Query().Get("salary")
	yearStr := r.URL.Query().Get("year")

	if salaryStr == "" || yearStr == "" {
		log.Printf("missing parameters salary and/or year: %v", salaryStr)
		http.Error(w, "missing parameters salary and/or year", http.StatusBadRequest)
		return 0, 0, errors.New("missing parameters salary and/or year")
	}

	year, err = strconv.Atoi(yearStr)
	if err != nil {
		log.Printf("invalid year: %v", err)
		http.Error(w, "invalid year ", http.StatusBadRequest)
		return 0, 0, errors.New("invalid year")
	}

	if !validYear(year) {
		log.Printf("The year: %v is not between 2019 to 2022, thus it is not support.", year)
		http.Error(w, "year is out of supported range (2019 to 2022)", http.StatusBadRequest)
		return 0, 0, errors.New("invalid year")
	}

	income, err = strconv.ParseFloat(salaryStr, 64)
	if err != nil || income < 0 || math.IsNaN(income) {
		log.Printf("invalid salary: %v", income)
		http.Error(w, "invalid salary", http.StatusBadRequest)
		return 0, 0, errors.New("invalid salary")
	}

	return
}

func validYear(year int) bool {
	if year > 2022 || year < 2019 {
		return false
	}
	return true
}
