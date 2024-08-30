package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Reservation struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
