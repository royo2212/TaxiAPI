package postgres

import (
	"context"
	"database/sql"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type RidePostgres struct {
	db *sql.DB
}

func NewRidePostgres(db *sql.DB) *RidePostgres {
	return &RidePostgres{db: db}
}

func (s *RidePostgres) SaveRide(ctx context.Context, r *entity.Ride) error {
	query := `INSERT INTO rides (passenger_id, driver_id, origin, destination, status)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING ride_id`

	var driverID interface{}
	if r.DriverID != nil {
		driverID = *r.DriverID
	} else {
		driverID = nil
	}

	err := s.db.QueryRowContext(ctx, query, r.PassengerID, driverID, r.Origin, r.Destination, r.Status).
		Scan(&r.RideID)
	if err != nil {
		return err
	}

	return nil
}

func (s *RidePostgres) FindRideByID(ctx context.Context, id int) (*entity.Ride, error) {
	query := `SELECT ride_id, passenger_id, driver_id, origin, destination, status
	          FROM rides
	          WHERE ride_id = $1`

	r := &entity.Ride{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&r.RideID, &r.PassengerID, &r.DriverID, &r.Origin, &r.Destination, &r.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrRideNotFound
		}
		return nil, err
	}
	return r, nil
}

func (s *RidePostgres) UpdateRideStatus(ctx context.Context, rideID int, status entity.Status) error {
	query := `UPDATE rides SET status = $1 WHERE ride_id = $2`

	result, err := s.db.ExecContext(ctx, query, status, rideID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return customErrors.ErrRideNotFound
	}

	return nil
}

func (s *RidePostgres) AssignDriverToRide(ctx context.Context, rideID, driverID int) error {
	query := `UPDATE rides SET driver_id = $1 WHERE ride_id = $2`

	result, err := s.db.ExecContext(ctx, query, driverID, rideID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return customErrors.ErrRideNotFound
	}

	return nil
}

func (s *RidePostgres) GetAllRides(ctx context.Context) ([]*entity.Ride, error) {
	query := `SELECT ride_id, passenger_id, driver_id, origin, destination, status
	          FROM rides`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*entity.Ride
	for rows.Next() {
		r := &entity.Ride{}
		if err := rows.Scan(&r.RideID, &r.PassengerID, &r.DriverID, &r.Origin, &r.Destination, &r.Status); err != nil {
			return nil, err
		}
		rides = append(rides, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rides, nil
}

func (s *RidePostgres) FindActiveRideByDriver(ctx context.Context, driverID int) (*entity.Ride, error) {
	query := `SELECT ride_id, passenger_id, driver_id, origin, destination, status
	          FROM rides
	          WHERE driver_id = $1 AND status IN ('pending', 'accepted')`

	r := &entity.Ride{}
	err := s.db.QueryRowContext(ctx, query, driverID).Scan(&r.RideID, &r.PassengerID, &r.DriverID, &r.Origin, &r.Destination, &r.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrRideNotFound
		}
		return nil, err
	}
	return r, nil
}

