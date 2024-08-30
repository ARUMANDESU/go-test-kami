package domain

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateReservationCreateDTO(value any) error {
	dto, ok := value.(ReservationCreateDTO)
	if !ok {
		return errors.New("invalid type")
	}

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.RoomID, validation.Required, is.UUID),
		validation.Field(&dto.StartTime, validation.Required),
		validation.Field(&dto.EndTime, validation.Required),
	)
	if err != nil {
		return err
	}

	if dto.EndTime.Before(dto.StartTime) {
		return errors.New("end time is before start time")
	}

	if dto.StartTime.Equal(dto.EndTime) {
		return errors.New("start time is equal to end time")
	}

	return nil
}
