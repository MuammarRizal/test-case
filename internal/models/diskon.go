package models

// DiskonRequest represents the request data for discount calculation
type DiskonRequest struct {
	Discounts          []DiskonItem `json:"discounts"`
	TotalSebelumDiskon float64      `json:"total_sebelum_diskon"`
}

// DiskonItem represents a single discount item
type DiskonItem struct {
	Diskon string `json:"diskon"`
}

// DiskonResponse represents the basic response for discount calculation
type DiskonResponse struct {
	TotalDiskon             float64 `json:"total_diskon"`
	TotalHargaSetelahDiskon float64 `json:"total_harga_setelah_diskon"`
}

// DiskonDetailResponse represents the detailed response for discount calculation
type DiskonDetailResponse struct {
	DiskonResponse
	TotalSebelumDiskon  float64             `json:"total_sebelum_diskon"`
	JumlahTingkatDiskon int                 `json:"jumlah_tingkat_diskon"`
	DetailPerhitungan   []DetailPerhitungan `json:"detail_perhitungan,omitempty"`
}

// DetailPerhitungan represents the step-by-step calculation details
type DetailPerhitungan struct {
	Step          int     `json:"step"`
	PersenDiskon  float64 `json:"persen_diskon"`
	HargaSebelum  float64 `json:"harga_sebelum"`
	NominalDiskon float64 `json:"nominal_diskon"`
	HargaSetelah  float64 `json:"harga_setelah"`
}