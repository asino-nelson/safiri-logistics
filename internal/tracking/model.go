package tracking

import "time"

type Event struct {
	ID             string    `json:"id"`
	LoadID         string    `json:"load_id"`
	DriverID       string    `json:"driver_id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	SpeedKPH       float64   `json:"speed_kph"`
	HeadingDegrees float64   `json:"heading_degrees"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type CreateEventRequest struct {
	Latitude       float64 `json:"latitude" binding:"required"`
	Longitude      float64 `json:"longitude" binding:"required"`
	SpeedKPH       float64 `json:"speed_kph"`
	HeadingDegrees float64 `json:"heading_degrees"`
}
