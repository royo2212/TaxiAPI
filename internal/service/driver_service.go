package service

import (
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type DriverStore interface {
	RegisterDriver(d *entity.Driver) (*entity.Driver, error)
	GetDriverByID(id int) (*entity.Driver, error)
	GetAllDrivers() ([]*entity.Driver, error)
	DeleteDriver(id int) error
	FindByPhoneNumber(phone int) (*entity.Driver, error)
}
type DriverService struct {
	store DriverStore
}

func NewDriverService(store DriverStore) *DriverService {
	return &DriverService{
		store: store,
	}
}
func (s *DriverService) RegisterDriver(d *entity.Driver) (*entity.Driver, error) {
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
	existing, _ := s.store.FindByPhoneNumber(d.PhoneNumber)
	if existing != nil {
		return nil, customErrors.ErrPhoneNumberExists
	}
	registeredDriver, err := s.store.RegisterDriver(d)
	if err != nil {
		return nil, err
	}
	return registeredDriver, nil
}

func (s *DriverService) GetDriverByID(id int) (*entity.Driver, error) {
	if id == 0 {
		return nil, customErrors.ErrDriverNotFound
	}
	driver, err := s.store.GetDriverByID(id)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (s *DriverService) GetAllDrivers() ([]*entity.Driver, error) {
	drivers, err := s.store.GetAllDrivers()
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (s *DriverService) DeleteDriver(id int) error {
	if id == 0 {
		return customErrors.ErrDriverNotFound
	}
	return s.store.DeleteDriver(id)
}
