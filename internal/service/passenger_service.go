package service

import (
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type PassengerStore interface {
	RegisterPassenger(p *entity.Passenger) (*entity.Passenger, error)
	GetPassengerByID(id int) (*entity.Passenger, error)
	GetAllPassengers() ([]*entity.Passenger, error)
	DeletePassenger(id int) error
	FindByPhoneNumber(phone int) (*entity.Passenger, error)
}

type PassengerService struct {
	store PassengerStore
}

func NewPassengerService(store PassengerStore) *PassengerService {
	return &PassengerService{
		store: store,
	}
}

func (s *PassengerService) RegisterPassenger(p *entity.Passenger) (*entity.Passenger, error) {
	if p == nil {
		return nil, customErrors.ErrPassengerDataRequired
	}
	if p.FirstName == "" {
		return nil, customErrors.ErrFirstName
	}
	if p.LastName == "" {
		return nil, customErrors.ErrLastName
	}
	if p.PhoneNumber == 0 {
		return nil, customErrors.ErrPhoneNumber
	}
	existing, _ := s.store.FindByPhoneNumber(p.PhoneNumber)
	if existing != nil {
		return nil, customErrors.ErrPhoneNumberExists
	}
	registeredPassenger, err := s.store.RegisterPassenger(p)
	if err != nil {
		return nil, err
	}
	return registeredPassenger, nil
}

func (s *PassengerService) GetPassengerByID(id int) (*entity.Passenger, error) {
	if id == 0 {
		return nil, customErrors.ErrPassengerNotFound
	}
	passenger, err := s.store.GetPassengerByID(id)
	if err != nil {
		return nil, err
	}
	return passenger, nil
}

func (s *PassengerService) GetAllPassengers() ([]*entity.Passenger, error) {
	passengers, err := s.store.GetAllPassengers()
	if err != nil {
		return nil, err
	}
	return passengers, nil
}

func (s *PassengerService) DeletePassenger(id int) error {
	if id == 0 {
		return customErrors.ErrPassengerNotFound
	}
	return s.store.DeletePassenger(id)
}
