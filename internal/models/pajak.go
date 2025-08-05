package models

type PajakRequest struct {
	Total       float64 `json:"total"`
	PersenPajak float64 `json:"persen_pajak"`
}

type PajakResponse struct {
	NetSales float64 `json:"net_sales"`
	PajakRp  float64 `json:"pajak_rp"`
}