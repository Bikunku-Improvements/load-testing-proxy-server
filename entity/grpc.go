package entity

import "time"

type BusLocationGRPC struct {
	Long      string    `json:"long"`
	Lat       string    `json:"lat"`
	CreatedAt time.Time `json:"created_at"`
	BusID     string    `json:"bus_id"`
}
