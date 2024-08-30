package domain

type ReservationCreateDTO struct {
	RoomID    string `json:"room_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
