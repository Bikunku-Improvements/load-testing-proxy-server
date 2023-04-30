package ws_client

import "time"

type BusLocationResponse struct {
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
