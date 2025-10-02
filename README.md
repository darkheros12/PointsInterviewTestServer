# Points Interview Test Server

This project is a Golang-based Tax Calculator API that calculates the total income tax owed for a given salary and year, using official marginal tax brackets.

## Features
- Tax calculation
- RESTful API endpoints
- In-memory caching for tax brackets
- Unit
- Automated CI with GitHub Actions

## Project Structure
Separation of concerns to improve maintainability and scalability
```
cmd/                # Application entry point
internal/
  client/           # HTTP client for tax brackets API (dockerized Mock API)
  controllers/      # HTTP handlers
  models/           # Data models
  services/         # Business logic
  utils/            # Utility functions (e.g., cache)
test/               # Integration and extra tests
go.mod, go.sum      # Go modules
```

## Requirements
- Go 1.25 or newer

## Running the Server
```
go run ./cmd/main.go
```

## Running Tests
```
go test ./...
```

## API Example Endpoint
```
GET {baseURL}/tax-calculator/tax-year?year=2022&salary=100000
```
Returns JSON with total tax, bands, and effective rate.

## API Example Request
```
curl "http://localhost:8080/tax-calculator?year=2022&salary=100000"
```

## API Example Response
```json
{
  "year": 2022,
  "income": 100000,
  "total_tax": 17739.17,
  "effective_rate": 0.18,
  "bands": [
    {
      "min": 0,
      "max": 50197,
      "rate": 0.15,
      "taxable": 50197,
      "tax": 7529.55
    },
    {
      "min": 50197,
      "max": 100392,
      "rate": 0.205,
      "taxable": 49803,
      "tax": 10209.62
    },
    {
      "min": 100392,
      "max": 155625,
      "rate": 0.26,
      "taxable": 0,
      "tax": 0
    },
    {
      "min": 155625,
      "max": 221708,
      "rate": 0.29,
      "taxable": 0,
      "tax": 0
    },
    {
      "min": 221708,
      "rate": 0.33,
      "taxable": 0,
      "tax": 0
    }
  ]
}
```

## Continuous Integration
GitHub Actions is configured to run all tests on Ubuntu, Windows, and macOS for every push and pull request.

