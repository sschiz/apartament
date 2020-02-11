package models

type House struct {
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
	Corpus   string `json:"corpus"`
	Floors   int    `json:"floors"`
}
