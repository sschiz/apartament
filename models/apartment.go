package models

type Apartment struct {
	Floor            int               `json:"floor"`
	Rooms            int               `json:"rooms"`
	Area             float64           `json:"area"`
	Rent             float64           `json:"rent"`
	House            *House            `json:"house"`
	ApartmentComplex *ApartmentComplex `json:"apartment_complex"`
}
