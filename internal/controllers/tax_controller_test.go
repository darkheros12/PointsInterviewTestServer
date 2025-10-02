package controllers

import (
	"net/http"
	"net/url"
	"testing"
)

type fakeWriter struct {
	headers http.Header
	status  int
	body    []byte
}

func (f *fakeWriter) Header() http.Header {
	if f.headers == nil {
		f.headers = make(http.Header)
	}
	return f.headers
}
func (f *fakeWriter) Write(b []byte) (int, error) {
	f.body = append(f.body, b...)
	return len(b), nil
}
func (f *fakeWriter) WriteHeader(statusCode int) {
	f.status = statusCode
}

func TestRequestValidation(t *testing.T) {

	tests := []struct {
		name       string
		params     map[string]string
		wantYear   int
		wantIncome float64
		wantErr    bool
		wantStatus int
	}{
		{
			"valid",
			map[string]string{"year": "2020", "salary": "12345.67"},
			2020,
			12345.67,
			false,
			0,
		},
		{
			"missing year",
			map[string]string{"salary": "1000"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
		{
			"missing salary",
			map[string]string{"year": "2021"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
		{
			"invalid year format",
			map[string]string{"year": "abcd", "salary": "1000"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
		{
			"invalid salary format",
			map[string]string{"year": "2021", "salary": "abc"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
		{
			"negative salary",
			map[string]string{"year": "2021", "salary": "-1"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
		{
			"year out of range",
			map[string]string{"year": "2018", "salary": "1000"},
			0,
			0,
			true,
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := &fakeWriter{}
			u := url.Values{}
			for key, value := range tt.params {
				u.Set(key, value)
			}
			r := &http.Request{URL: &url.URL{RawQuery: u.Encode()}}
			year, income, err := requestValidation(fw, r)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if year != tt.wantYear {
				t.Errorf("got year %d, want %d", year, tt.wantYear)
			}
			if income != tt.wantIncome {
				t.Errorf("got income %f, want %f", income, tt.wantIncome)
			}
			if tt.wantErr && fw.status != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, fw.status)
			}
		})
	}
}

func TestValidYear(t *testing.T) {
	cases := []struct {
		year  int
		valid bool
	}{
		{2019, true},
		{2020, true},
		{2022, true},
		{2018, false},
		{2023, false},
	}
	for _, c := range cases {
		if got := validYear(c.year); got != c.valid {
			t.Errorf("validYear(%d) = %v, want %v", c.year, got, c.valid)
		}
	}
}
