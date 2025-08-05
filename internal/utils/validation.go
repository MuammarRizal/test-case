package utils

import (
	"fmt"
	"net/http"
)

// ValidateMethod checks if the request method is allowed
func ValidateMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}
	RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	return false
}

// ConvertToString converts interface{} to string
func ConvertToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}