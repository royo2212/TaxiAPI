package postgres

import (
	"context"
	"database/sql"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type PassengerPostgres struct {
	db *sql.DB
}

func NewPassengerPostgres(db *sql.DB) *PassengerPostgres {
	return &PassengerPostgres{
		db: db,
	}
}
func (s *PassengerPostgres) RegisterPassenger(ctx context.Context, p *entity.Passenger) (*entity.Passenger, error) {
	query := `INSERT INTO passengers (first_name, last_name, phone_number)
	          VALUES ($1, $2, $3)
	          RETURNING id`
	err := s.db.QueryRowContext(ctx, query, p.FirstName, p.LastName, p.PhoneNumber).Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (s *PassengerPostgres) GetPassengerByID(ctx context.Context, id int) (*entity.Passenger, error) {
	query := `SELECT id, first_name, last_name, phone_number
	          FROM passengers
	          WHERE id = $1`

	p := &entity.Passenger{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.FirstName, &p.LastName, &p.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrPassengerNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *PassengerPostgres) GetAllPassengers(ctx context.Context) ([]*entity.Passenger, error) {
	query := `SELECT id, first_name, last_name, phone_number FROM passengers`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []*entity.Passenger
	for rows.Next() {
		p := &entity.Passenger{}
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.PhoneNumber); err != nil {
			return nil, err
		}
		passengers = append(passengers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return passengers, nil
}
func (s *PassengerPostgres) DeletePassenger(ctx context.Context, id int) error {
	query := `DELETE FROM passengers WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return customErrors.ErrPassengerNotFound
	}

	return nil
}

func (s *PassengerPostgres) FindByPhoneNumber(ctx context.Context, phone int) (*entity.Passenger, error) {
	query := `SELECT id, first_name, last_name, phone_number
	          FROM passengers
	          WHERE phone_number = $1`

	p := &entity.Passenger{}
	err := s.db.QueryRowContext(ctx, query, phone).Scan(&p.ID, &p.FirstName, &p.LastName, &p.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrPassengerNotFound
		}
		return nil, err
	}
	return p, nil
}
