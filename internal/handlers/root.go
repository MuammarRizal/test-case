package handlers

import (
	"encoding/json"
	"net/http"
)

// RootHandler handles the root endpoint and provides API documentation
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "API Server - Sales, Tax & Discount Calculator",
		"endpoints": map[string]interface{}{
			"POST /penjualan":     "Save sales transaction",
			"GET /hitung-pajak":   "Calculate tax with query parameters",
			"POST /hitung-pajak":  "Calculate tax with JSON body",
			"POST /hitung-diskon": "Calculate tiered discount",
		},
		"examples": map[string]interface{}{
			"tax_calculation": map[string]string{
				"GET":  "curl 'http://localhost:8080/hitung-pajak?total=22000&persen_pajak=10'",
				"POST": "curl -X POST http://localhost:8080/hitung-pajak -H 'Content-Type: application/json' -d '{\"total\":22000,\"persen_pajak\":10}'",
			},
			"discount_calculation": "curl -X POST http://localhost:8080/hitung-diskon -H 'Content-Type: application/json' -d '{\"discounts\":[{\"diskon\":\"20\"},{\"diskon\":\"10\"}],\"total_sebelum_diskon\":100000}'",
		},
	})
}