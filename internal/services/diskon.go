package services

import (
	"fmt"
	"strconv"

	"test-case/internal/models"
	"test-case/internal/utils"
)

func ValidateDiskonInput(request models.DiskonRequest) error {
	if request.TotalSebelumDiskon <= 0 {
		return fmt.Errorf("total_sebelum_diskon harus lebih besar dari 0")
	}

	if len(request.Discounts) == 0 {
		return fmt.Errorf("discounts tidak boleh kosong")
	}

	for i, discount := range request.Discounts {
		if discount.Diskon == "" {
			return fmt.Errorf("diskon pada index %d tidak boleh kosong", i)
		}

		diskonPersen, err := strconv.ParseFloat(discount.Diskon, 64)
		if err != nil {
			return fmt.Errorf("diskon pada index %d harus berupa angka valid", i)
		}

		if diskonPersen < 0 {
			return fmt.Errorf("diskon pada index %d tidak boleh negatif", i)
		}

		if diskonPersen > 100 {
			return fmt.Errorf("diskon pada index %d tidak boleh lebih dari 100%%", i)
		}
	}

	return nil
}

func HitungDiskonBertingkat(totalSebelumDiskon float64, discounts []models.DiskonItem) models.DiskonResponse {
	hargaSekarang := totalSebelumDiskon

	for _, discount := range discounts {
		diskonPersen, _ := strconv.ParseFloat(discount.Diskon, 64)
		nominalDiskon := hargaSekarang * (diskonPersen / 100)
		hargaSekarang = hargaSekarang - nominalDiskon
	}

	totalDiskon := totalSebelumDiskon - hargaSekarang

	return models.DiskonResponse{
		TotalDiskon:             utils.RoundToTwoDecimals(totalDiskon),
		TotalHargaSetelahDiskon: utils.RoundToTwoDecimals(hargaSekarang),
	}
}

func HitungDiskonBertingkatDetail(totalSebelumDiskon float64, discounts []models.DiskonItem) models.DiskonDetailResponse {
	hargaSekarang := totalSebelumDiskon
	var detailPerhitungan []models.DetailPerhitungan

	for i, discount := range discounts {
		diskonPersen, _ := strconv.ParseFloat(discount.Diskon, 64)
		hargaSebelum := hargaSekarang
		nominalDiskon := hargaSekarang * (diskonPersen / 100)
		hargaSekarang = hargaSekarang - nominalDiskon

		detailPerhitungan = append(detailPerhitungan, models.DetailPerhitungan{
			Step:          i + 1,
			PersenDiskon:  diskonPersen,
			HargaSebelum:  utils.RoundToTwoDecimals(hargaSebelum),
			NominalDiskon: utils.RoundToTwoDecimals(nominalDiskon),
			HargaSetelah:  utils.RoundToTwoDecimals(hargaSekarang),
		})
	}

	totalDiskon := totalSebelumDiskon - hargaSekarang

	return models.DiskonDetailResponse{
		DiskonResponse: models.DiskonResponse{
			TotalDiskon:             utils.RoundToTwoDecimals(totalDiskon),
			TotalHargaSetelahDiskon: utils.RoundToTwoDecimals(hargaSekarang),
		},
		TotalSebelumDiskon:  totalSebelumDiskon,
		JumlahTingkatDiskon: len(discounts),
		DetailPerhitungan:   detailPerhitungan,
	}
}