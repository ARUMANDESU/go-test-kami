package domain

import (
	"encoding/json"
	"time"
)

type ReservationCreateDTO struct {
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// MarshalJSON converts the ReservationCreateDTO to JSON
// with the StartTime and EndTime fields formatted as RFC3339.
func (r ReservationCreateDTO) MarshalJSON() ([]byte, error) {
	type Alias ReservationCreateDTO
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

// UnmarshalJSON converts the JSON data to a ReservationCreateDTO.
// It expects the StartTime and EndTime fields to be formatted as RFC3339.
func (r *ReservationCreateDTO) UnmarshalJSON(data []byte) error {
	type Alias ReservationCreateDTO
	aux := &struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		*Alias
	}{
		Alias:     (*Alias)(r),
		StartTime: r.StartTime.Format(time.RFC3339),
		EndTime:   r.EndTime.Format(time.RFC3339),
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
