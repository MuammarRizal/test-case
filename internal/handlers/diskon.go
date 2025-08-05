package handlers

import (
	"encoding/json"
	"net/http"

	"test-case/internal/models"
	"test-case/internal/services"
	"test-case/internal/utils"
)

func DiskonHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	var request models.DiskonRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	if err := services.ValidateDiskonInput(request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	showDetail := r.URL.Query().Get("detail") == "true"

	if showDetail {
		result := services.HitungDiskonBertingkatDetail(request.TotalSebelumDiskon, request.Discounts)
		utils.RespondSuccess(w, "Perhitungan diskon bertingkat berhasil", result)
	} else {
		result := services.HitungDiskonBertingkat(request.TotalSebelumDiskon, request.Discounts)
		utils.RespondSuccess(w, "Perhitungan diskon bertingkat berhasil", result)
	}
}