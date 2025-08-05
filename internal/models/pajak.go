package models

// PajakRequest represents the request data for tax calculation
type PajakRequest struct {
	Total       float64 `json:"total"`
	PersenPajak float64 `json:"persen_pajak"`
}

// PajakResponse represents the response data for tax calculation
type PajakResponse struct {
	NetSales float64 `json:"net_sales"`
	PajakRp  float64 `json:"pajak_rp"`
}