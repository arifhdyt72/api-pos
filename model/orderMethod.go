package model

type OrderMethod struct {
	MavisModel
	Name      string `json:"name"`
	MarkOrder string `json:"mark_order"`
	Price     int64  `json:"price"`
}
