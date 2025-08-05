package services

import (
	"fmt"

	"test-case/internal/models"
	"test-case/internal/utils"
)

func ValidatePajakInput(request models.PajakRequest) error {
	if request.Total <= 0 {
		return fmt.Errorf("total harus lebih besar dari 0")
	}

	if request.PersenPajak < 0 {
		return fmt.Errorf("persen_pajak tidak boleh negatif")
	}

	if request.PersenPajak > 100 {
		return fmt.Errorf("persen_pajak tidak boleh lebih dari 100")
	}

	return nil
}

func HitungPajak(total float64, persenPajak float64) models.PajakResponse {
	netSales := total / (1 + (persenPajak / 100))
	pajakRp := total - netSales

	return models.PajakResponse{
		NetSales: utils.RoundToTwoDecimals(netSales),
		PajakRp:  utils.RoundToTwoDecimals(pajakRp),
	}
}