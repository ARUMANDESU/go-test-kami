package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Reservation struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// OverlapsWith checks if the reservation overlaps with another reservation.
//
//	other.startTime < r.endTime && r.startTime < other.endTime
//
// A reservation is considered to be in conflict with another reservation if:
//   - the start time of the reservation is before the end time of the other reservation
//   - the end time of the reservation is after the start time of the other reservation
func (r Reservation) OverlapsWith(other Reservation) bool {
	return other.StartTime.Before(r.EndTime) && r.StartTime.Before(other.EndTime)
}
