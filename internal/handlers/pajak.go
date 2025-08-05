package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"test-case/internal/models"
	"test-case/internal/services"
	"test-case/internal/utils"
)

// PajakHandler handles tax calculation requests (both GET and POST)
func PajakHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPajak(w, r)
	case http.MethodPost:
		handlePostPajak(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed")
	}
}

func handleGetPajak(w http.ResponseWriter, r *http.Request) {
	totalStr := r.URL.Query().Get("total")
	persenPajakStr := r.URL.Query().Get("persen_pajak")

	if totalStr == "" || persenPajakStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "Parameter 'total' dan 'persen_pajak' harus diisi")
		return
	}

	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Parameter 'total' harus berupa angka")
		return
	}

	persenPajak, err := strconv.ParseFloat(persenPajakStr, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Parameter 'persen_pajak' harus berupa angka")
		return
	}

	request := models.PajakRequest{Total: total, PersenPajak: persenPajak}
	if err := services.ValidatePajakInput(request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	result := services.HitungPajak(request.Total, request.PersenPajak)
	utils.RespondSuccess(w, "Perhitungan pajak berhasil", result)
}

func handlePostPajak(w http.ResponseWriter, r *http.Request) {
	var request models.PajakRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	if err := services.ValidatePajakInput(request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	result := services.HitungPajak(request.Total, request.PersenPajak)
	utils.RespondSuccess(w, "Perhitungan pajak berhasil", result)
}