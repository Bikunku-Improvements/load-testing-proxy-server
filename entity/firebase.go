package entity

import "time"

type Bus struct {
	ID       int    `json:"id"`
	Number   int    `json:"number"`
	Plate    string `json:"plate"`
	Status   string `json:"status"`
	Route    string `json:"route"`
	IsActive bool   `json:"is_active"`
}

type BusLocationFirebase struct {
	BusID     int       `json:"bus_id"`
	Heading   int       `json:"heading"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     int       `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

type BusLocationFirebaseResponse struct {
	Number    int       `json:"number"`
	Plate     string    `json:"plate"`
	Status    string    `json:"status"`
	Route     string    `json:"route"`
	IsActive  bool      `json:"is_active"`
	Heading   int       `json:"heading"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     int       `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}
