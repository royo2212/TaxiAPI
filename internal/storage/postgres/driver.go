package postgres
import (
	"context"
	"database/sql"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type DriverPostgres struct {
	db *sql.DB
}

func NewDriverPostgres(db *sql.DB) *DriverPostgres {
	return &DriverPostgres{db: db}
}
func (s *DriverPostgres) RegisterDriver(ctx context.Context, d *entity.Driver) (*entity.Driver, error) {
	query := `INSERT INTO drivers (first_name, last_name, phone_number, car_type, license_plate, is_available)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id`
	err := s.db.QueryRowContext(ctx, query, d.FirstName, d.LastName, d.PhoneNumber, d.CarType, d.LicensePlate, d.IsAvailable).Scan(&d.ID)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DriverPostgres) GetDriverByID(ctx context.Context, id int) (*entity.Driver, error) {
	query := `SELECT id, first_name, last_name, phone_number, car_type, license_plate, is_available
	          FROM drivers
	          WHERE id = $1`

	d := &entity.Driver{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.FirstName, &d.LastName, &d.PhoneNumber, &d.CarType, &d.LicensePlate, &d.IsAvailable)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrDriverNotFound
		}
		return nil, err
	}
	return d, nil
}

func (s *DriverPostgres) GetAllDrivers(ctx context.Context) ([]*entity.Driver, error) {
	query := `SELECT id, first_name, last_name, phone_number, car_type, license_plate, is_available
	          FROM drivers`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []*entity.Driver
	for rows.Next() {
		d := &entity.Driver{}
		if err := rows.Scan(&d.ID, &d.FirstName, &d.LastName, &d.PhoneNumber, &d.CarType, &d.LicensePlate, &d.IsAvailable); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (s *DriverPostgres) DeleteDriver(ctx context.Context, id int) error {
	query := `DELETE FROM drivers WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customErrors.ErrDriverNotFound
	}
	return nil
}

func (s *DriverPostgres) FindByPhoneNumber(ctx context.Context, phone int) (*entity.Driver, error) {
	query := `SELECT id, first_name, last_name, phone_number, car_type, license_plate, is_available
	          FROM drivers
	          WHERE phone_number = $1`

	d := &entity.Driver{}
	err := s.db.QueryRowContext(ctx, query, phone).Scan(&d.ID, &d.FirstName, &d.LastName, &d.PhoneNumber, &d.CarType, &d.LicensePlate, &d.IsAvailable)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrDriverNotFound
		}
		return nil, err
	}
	return d, nil
}
