package entity

import "time"

type BusLocationWSLegacyResponse struct {
	Id            int       `json:"id"`
	Number        int       `json:"number"`
	Plate         string    `json:"plate"`
	Status        string    `json:"status"`
	Route         string    `json:"route"`
	IsActive      bool      `json:"isActive"`
	Long          int       `json:"long"`
	Lat           int       `json:"lat"`
	Speed         int       `json:"speed"`
	Heading       int       `json:"heading"`
	Timestamp     time.Time `json:"Timestamp"`
	IsNewLocation bool      `json:"IsNewLocation"`
}

type SendLocationRequest struct {
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
	Speed   float64 `json:"speed"`
	Heading float64 `json:"heading"`
}
