package domain

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateReservationCreateDTO(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		name    string
		value   any
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid DTO",
			value: ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow.Add(1 * time.Hour),
				EndTime:   timeNow.Add(2 * time.Hour),
			},
			wantErr: false,
		},
		{
			name:    "invalid type",
			value:   "invalid",
			wantErr: true,
			errMsg:  "invalid type",
		},
		{
			name: "missing RoomID",
			value: ReservationCreateDTO{
				StartTime: timeNow.Add(1 * time.Hour),
				EndTime:   timeNow.Add(2 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "invalid RoomID format",
			value: ReservationCreateDTO{
				RoomID:    "invalid-uuid",
				StartTime: timeNow.Add(1 * time.Hour),
				EndTime:   timeNow.Add(2 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "missing StartTime",
			value: ReservationCreateDTO{
				RoomID:  uuid.Must(uuid.NewV7()).String(),
				EndTime: timeNow.Add(2 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "missing EndTime",
			value: ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow.Add(1 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "StartTime is before EndTime",
			value: ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow.Add(2 * time.Hour),
				EndTime:   timeNow.Add(1 * time.Hour),
			},
			wantErr: true,
			errMsg:  "end time is before start time",
		},
		{
			name: "StartTime is equal to EndTime",
			value: ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow.Add(1 * time.Hour),
				EndTime:   timeNow.Add(1 * time.Hour),
			},
			wantErr: true,
			errMsg:  "start time is equal to end time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReservationCreateDTO(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
