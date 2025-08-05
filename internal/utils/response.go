package utils

import (
	"encoding/json"
	"math"
	"net/http"

	"test-case/internal/models"
)

// RespondJSON sends a JSON response with the given status code
func RespondJSON(w http.ResponseWriter, statusCode int, response models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, models.APIResponse{
		Status:  "error",
		Message: message,
	})
}

// RespondSuccess sends a success response
func RespondSuccess(w http.ResponseWriter, message string, data interface{}) {
	RespondJSON(w, http.StatusOK, models.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RoundToTwoDecimals rounds a float64 value to two decimal places
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}