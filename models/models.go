package models

type Customer struct {
	Material     string  `json:"name"`
	AddressLat   float64 `json:"latitude"`
	AddressLong  float64 `json:"longitude"`
	SquareMeters float64 `json:"squaremeters"`
	PhoneNumber  int64   `json:"phonenumber"`
}

type Partner struct {
	PartnerID       int64   `json:"id"`
	PartnerName     string  `json:"name"`
	OperatingRadius int64   `json:"radius"`
	Rating          float64 `json:"rating"`
	AddressLat      float64 `json:"lattitude"`
	AddressLong     float64 `json:"longitude"`
}

type Result struct {
	ID              int64   `json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	OperatingRadius int64   `json:"OperatingRadius,omitempty"`
	Rating          float64 `json:"rating,omitempty"`
	AddressLat      float64 `json:"lattitude"`
	AddressLong     float64 `json:"longitude"`
	Distance        float64 `json:"distance,omitempty"`
}
