package models

type ApartmentComplex struct {
	Name       string `json:"name"`
	Apartments [2]int `json:"min_max_apartments"`
}
