package models

type ApartmentComplex struct {
	Name   string
	Houses []*House
	Floors [2]int
}
