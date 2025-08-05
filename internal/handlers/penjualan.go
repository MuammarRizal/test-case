package handlers

import (
	"encoding/json"
	"net/http"

	"test-case/internal/models"
	"test-case/internal/services"
	"test-case/internal/utils"
)

// PenjualanHandler handles sales transaction requests
func PenjualanHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	var data models.PenjualanRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	if err := services.ValidatePenjualanData(data); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	penjualanID, err := services.SavePenjualan(data)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to save penjualan: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, models.APIResponse{
		Status:  "success",
		Message: "Penjualan berhasil disimpan",
		Data:    map[string]interface{}{"penjualan_id": penjualanID},
	})
}