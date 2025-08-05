package services

import (
	"database/sql"
	"fmt"
	"time"

	"test-case/internal/config"
	"test-case/internal/models"
	"test-case/internal/utils"
)

func ValidatePenjualanData(data models.PenjualanRequest) error {
	if data.NamaPelanggan == "" {
		return fmt.Errorf("nama_pelanggan tidak boleh kosong")
	}

	if data.Tanggal == "" {
		return fmt.Errorf("tanggal tidak boleh kosong")
	}

	if _, err := time.Parse("2006-01-02", data.Tanggal); err != nil {
		return fmt.Errorf("format tanggal harus YYYY-MM-DD")
	}

	if data.Jam == "" {
		return fmt.Errorf("jam tidak boleh kosong")
	}

	if _, err := time.Parse("15:04", data.Jam); err != nil {
		return fmt.Errorf("format jam harus HH:MM")
	}

	if data.Total <= 0 {
		return fmt.Errorf("total harus lebih besar dari 0")
	}

	if data.BayarTunai < data.Total {
		return fmt.Errorf("bayar_tunai tidak boleh kurang dari total")
	}

	if len(data.Items) == 0 {
		return fmt.Errorf("items tidak boleh kosong")
	}

	var sumSubTotal float64
	for i, item := range data.Items {
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity item ke-%d harus lebih besar dari 0", i+1)
		}
		if item.Harga <= 0 {
			return fmt.Errorf("harga item ke-%d harus lebih besar dari 0", i+1)
		}
		if item.SubTotal <= 0 {
			return fmt.Errorf("sub_total item ke-%d harus lebih besar dari 0", i+1)
		}
		sumSubTotal += item.SubTotal
	}

	if sumSubTotal != data.Total {
		return fmt.Errorf("total tidak sesuai dengan jumlah sub_total items")
	}

	return nil
}

func SavePenjualan(data models.PenjualanRequest) (int64, error) {
	db := config.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	penjualanID, err := insertPenjualan(tx, data)
	if err != nil {
		return 0, err
	}

	err = insertPenjualanItems(tx, penjualanID, data.Items)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return penjualanID, nil
}

func insertPenjualan(tx *sql.Tx, data models.PenjualanRequest) (int64, error) {
	query := `
		INSERT INTO penjualan (nama_pelanggan, tanggal, jam, total, bayar_tunai, kembali)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query,
		data.NamaPelanggan, data.Tanggal, data.Jam,
		data.Total, data.BayarTunai, data.Kembali,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func insertPenjualanItems(tx *sql.Tx, penjualanID int64, items []models.PenjualanItem) error {
	query := `
		INSERT INTO penjualan_item (penjualan_id, item_id, quantity, harga, sub_total)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		itemID := utils.ConvertToString(item.ItemID)
		_, err := stmt.Exec(penjualanID, itemID, item.Quantity, item.Harga, item.SubTotal)
		if err != nil {
			return err
		}
	}

	return nil
}