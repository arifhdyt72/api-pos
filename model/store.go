package model

type Store struct {
	MavisModel
	Name    string `json:"name"`
	Address string `json:"address"`
}
