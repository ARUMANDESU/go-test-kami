package domain

import (
	"encoding/json"
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

// MarshalJSON converts the Reservation to JSON
// with the StartTime and EndTime fields formatted as RFC3339.
func (r Reservation) MarshalJSON() ([]byte, error) {
	type Alias Reservation
	return json.Marshal(&struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		*Alias
	}{
		StartTime: r.StartTime.Format(time.RFC3339),
		EndTime:   r.EndTime.Format(time.RFC3339),
		Alias:     (*Alias)(&r),
	})
}

// UnmarshalJSON converts the JSON data to a Reservation.
// It expects the StartTime and EndTime fields to be formatted as RFC3339.
func (r Reservation) UnmarshalJSON(data []byte) error {
	type Alias Reservation
	aux := &struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		*Alias
	}{
		Alias: (*Alias)(&r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	startTime, err := time.Parse(time.RFC3339, aux.StartTime)
	if err != nil {
		return err
	}
	r.StartTime = startTime

	endTime, err := time.Parse(time.RFC3339, aux.EndTime)
	if err != nil {
		return err
	}
	r.EndTime = endTime

	return nil
}
