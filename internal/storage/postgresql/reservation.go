package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetRoomReservations returns all reservations for a given room.
// If no reservations are found, the function returns an `ErrNotFound` error.
func (s Storage) GetRoomReservations(ctx context.Context, roomID string) ([]domain.Reservation, error) {
	const op = "storage.postgresql.GetRoomReservations"

	query := `
		SELECT  id, room_id, start_time, end_time
		FROM    reservations
		WHERE   room_id = $1
	`

	rows, err := s.Pool.Query(ctx, query, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %s: %w", op, "failed to query database", err)
	}
	defer rows.Close()

	var reservations []domain.Reservation

	for rows.Next() {
		var r domain.Reservation

		err := rows.Scan(&r.ID, &r.RoomID, &r.StartTime, &r.EndTime)
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", op, "failed to scan row", err)
		}

		reservations = append(reservations, r)
	}

	return reservations, nil
}

// ReserveRoom reserves a room for a given time range.
// It returns the reservation if the room is available, otherwise it returns an error.
// The function uses a serializable transaction to ensure that the reservation is atomic.
// If the reservation conflicts with another reservation, the function returns an `ErrResevationConflict` error.
func (s Storage) ReserveRoom(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error) {
	const op = "storage.postgresql.ReserveRoom"
	const maxRetries = 3

	var err error
	for i := 0; i < maxRetries; i++ {
		tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
		if err != nil {
			return domain.Reservation{}, fmt.Errorf("%s: %s: %w", op, "failed to begin transaction", err)
		}

		query := `
			SELECT COUNT(*)
			FROM   reservations
			WHERE  room_id = $1 AND (start_time, end_time) OVERLAPS ($2, $3)
		`

		var count int
		err = tx.QueryRow(ctx, query, reservation.RoomID, reservation.StartTime, reservation.EndTime).Scan(&count)
		if err != nil {
			tx.Rollback(ctx)
			return domain.Reservation{}, fmt.Errorf("%s: %s: %w", op, "failed to query database", err)
		}

		if count > 0 {
			tx.Rollback(ctx)
			return domain.Reservation{}, fmt.Errorf("%s: %w", op, storage.ErrResevationConflict)
		}

		query = `
			INSERT INTO reservations (id, room_id, start_time, end_time, created_at, updated_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`

		result, err := tx.Exec(ctx, query, reservation.ID, reservation.RoomID, reservation.StartTime, reservation.EndTime)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				switch pgErr.Code {
				case pgerrcode.SerializationFailure:
					tx.Rollback(ctx)
					// Retry on serialization failure
					continue
				default:
					tx.Rollback(ctx)
					return domain.Reservation{}, fmt.Errorf("%s: %s: %w", op, "failed to insert reservation", err)
				}
			}
		}

		// check if the reservation was inserted
		if result.RowsAffected() != 1 {
			tx.Rollback(ctx)
			return domain.Reservation{}, fmt.Errorf("%s: %s", op, "failed to insert reservation (no rows affected)")
		}

		err = tx.Commit(ctx)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				switch pgErr.Code {
				case pgerrcode.SerializationFailure:
					tx.Rollback(ctx)
					// Retry on serialization failure
					continue
				default:
					tx.Rollback(ctx)
					return domain.Reservation{}, fmt.Errorf("%s: %s: %w", op, "failed to commit transaction", err)
				}
			}
		}

		// If commit was successful, break the loop
		break
	}

	if err != nil {
		return domain.Reservation{}, fmt.Errorf("%s: %s: %w", op, "transaction failed after retries", err)
	}

	return reservation, nil
}
