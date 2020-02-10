package models

type Apartment struct {
	Floor            int
	Rooms            int
	Area             float64
	Rent             float64
	House            *House
	ApartmentComplex *ApartmentComplex
}
