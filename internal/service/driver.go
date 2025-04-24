package service

import (
	"context"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type DriverStore interface {
	RegisterDriver(ctx context.Context, d *entity.Driver) (*entity.Driver, error)
	GetDriverByID(ctx context.Context, id int) (*entity.Driver, error)
	GetAllDrivers(ctx context.Context) ([]*entity.Driver, error)
	DeleteDriver(ctx context.Context, id int) error
	FindByPhoneNumber(ctx context.Context, phone int) (*entity.Driver, error)
}

type DriverService struct {
	store DriverStore
}

func NewDriverService(store DriverStore) *DriverService {
	return &DriverService{
		store: store,
	}
}
func (s *DriverService) RegisterDriver(ctx context.Context, d *entity.Driver) (*entity.Driver, error) {
	if d == nil {
		return nil, customErrors.ErrDriverDataRequired
	}
	if d.FirstName == "" {
		return nil, customErrors.ErrFirstName
	}
	if d.LastName == "" {
		return nil, customErrors.ErrLastName
	}
	if d.CarType == "" {
		return nil, customErrors.ErrCarTypeRequired
	}
	if d.LicensePlate == 0 {
		return nil, customErrors.ErrLicensePlateRequired
	}
	if d.PhoneNumber == 0 {
		return nil, customErrors.ErrPhoneNumber
	}
	existing, _ := s.store.FindByPhoneNumber(ctx, d.PhoneNumber)
	if existing != nil {
		return nil, customErrors.ErrPhoneNumberExists
	}
	registeredDriver, err := s.store.RegisterDriver(ctx, d)
	if err != nil {
		return nil, err
	}
	return registeredDriver, nil
}

func (s *DriverService) GetDriverByID(ctx context.Context, id int) (*entity.Driver, error) {
	if id == 0 {
		return nil, customErrors.ErrDriverNotFound
	}
	driver, err := s.store.GetDriverByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (s *DriverService) GetAllDrivers(ctx context.Context) ([]*entity.Driver, error) {
	drivers, err := s.store.GetAllDrivers(ctx)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (s *DriverService) DeleteDriver(ctx context.Context, id int) error {
	if id == 0 {
		return customErrors.ErrDriverNotFound
	}
	return s.store.DeleteDriver(ctx, id)
}
