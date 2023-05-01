package app

import (
	"time"
)

type (
	RoomBooked struct {
		ReservationId string     `json:"reservation_id,omitempty"`
		RoomId        string     `json:"room_id,omitempty"`
		GuestName     string     `json:"guest_name,omitempty"`
		Price         int64      `json:"price,omitempty"`
		StartDate     *time.Time `json:"start_date,omitempty"`
		EndDate       *time.Time `json:"end_date,omitempty"`
	}

	BeerOrdered struct {
		RoomId string `json:"room_id,omitempty"`
		Count  int64  `json:"count,omitempty"`
	}
)
