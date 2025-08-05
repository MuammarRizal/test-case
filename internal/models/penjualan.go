package models

type PenjualanRequest struct {
	NamaPelanggan string          `json:"nama_pelanggan"`
	Tanggal       string          `json:"tanggal"`
	Jam           string          `json:"jam"`
	Total         float64         `json:"total"`
	BayarTunai    float64         `json:"bayar_tunai"`
	Kembali       float64         `json:"kembali"`
	Items         []PenjualanItem `json:"items"`
}

type PenjualanItem struct {
	ItemID   interface{} `json:"item_id"`
	Quantity float64     `json:"quantity"`
	Harga    float64     `json:"harga"`
	SubTotal float64     `json:"sub_total"`
}