package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PenjualanRequest struct {
	NamaPelanggan string           `json:"nama_pelanggan"`
	Tanggal       string           `json:"tanggal"`
	Jam           string           `json:"jam"`
	Total         float64          `json:"total"`
	BayarTunai    float64          `json:"bayar_tunai"`
	Kembali       float64          `json:"kembali"`
	Items         []PenjualanItem  `json:"items"`
}

type PenjualanItem struct {
	ItemID    interface{} `json:"item_id"`
	Quantity  float64     `json:"quantity"`
	Harga     float64     `json:"harga"`
	SubTotal  float64     `json:"sub_total"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/test_case?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	
	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
}

func postPenjualan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}
	
	var dataPenjualan PenjualanRequest
	if err := json.NewDecoder(r.Body).Decode(&dataPenjualan); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Invalid JSON format: " + err.Error(),
		})
		return
	}
	
	if err := validatePenjualanData(dataPenjualan); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Validation error: " + err.Error(),
		})
		return
	}
	
	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Failed to start transaction: " + err.Error(),
		})
		return
	}
	
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	penjualanID, err := insertPenjualan(tx, dataPenjualan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Failed to insert penjualan: " + err.Error(),
		})
		return
	}
	
	err = insertPenjualanItems(tx, penjualanID, dataPenjualan.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Failed to insert penjualan items: " + err.Error(),
		})
		return
	}
	
	if err = tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Failed to commit transaction: " + err.Error(),
		})
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Penjualan berhasil disimpan",
		Data: map[string]interface{}{
			"penjualan_id": penjualanID,
		},
	})
}

func validatePenjualanData(data PenjualanRequest) error {
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

func insertPenjualan(tx *sql.Tx, data PenjualanRequest) (int64, error) {
	query := `
		INSERT INTO penjualan (nama_pelanggan, tanggal, jam, total, bayar_tunai, kembali)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := tx.Exec(query, 
		data.NamaPelanggan,
		data.Tanggal,
		data.Jam,
		data.Total,
		data.BayarTunai,
		data.Kembali,
	)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

func insertPenjualanItems(tx *sql.Tx, penjualanID int64, items []PenjualanItem) error {
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
		var itemID string
		switch v := item.ItemID.(type) {
		case string:
			itemID = v
		case float64:
			itemID = fmt.Sprintf("%.0f", v)
		case int:
			itemID = fmt.Sprintf("%d", v)
		default:
			itemID = fmt.Sprintf("%v", v)
		}
		
		_, err := stmt.Exec(penjualanID, itemID, item.Quantity, item.Harga, item.SubTotal)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/penjualan", postPenjualan)
	
	// Start server
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}