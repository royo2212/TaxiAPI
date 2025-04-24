package service

import (
	"context"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type PassengerStore interface {
	RegisterPassenger(ctx context.Context, p *entity.Passenger) (*entity.Passenger, error)
	GetPassengerByID(ctx context.Context, id int) (*entity.Passenger, error)
	GetAllPassengers(ctx context.Context) ([]*entity.Passenger, error)
	DeletePassenger(ctx context.Context, id int) error
	FindByPhoneNumber(ctx context.Context, phone int) (*entity.Passenger, error)
}

type PassengerService struct {
	store PassengerStore
}

func NewPassengerService(store PassengerStore) *PassengerService {
	return &PassengerService{
		store: store,
	}
}

func (s *PassengerService) RegisterPassenger(ctx context.Context, p *entity.Passenger) (*entity.Passenger, error) {
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
	existing, _ := s.store.FindByPhoneNumber(ctx, p.PhoneNumber)
	if existing != nil {
		return nil, customErrors.ErrPhoneNumberExists
	}
	registeredPassenger, err := s.store.RegisterPassenger(ctx, p)
	if err != nil {
		return nil, err
	}
	return registeredPassenger, nil
}

func (s *PassengerService) GetPassengerByID(ctx context.Context, id int) (*entity.Passenger, error) {
	if id == 0 {
		return nil, customErrors.ErrPassengerNotFound
	}
	passenger, err := s.store.GetPassengerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return passenger, nil
}

func (s *PassengerService) GetAllPassengers(ctx context.Context) ([]*entity.Passenger, error) {
	passengers, err := s.store.GetAllPassengers(ctx)
	if err != nil {
		return nil, err
	}
	return passengers, nil
}

func (s *PassengerService) DeletePassenger(ctx context.Context, id int) error {
	if id == 0 {
		return customErrors.ErrPassengerNotFound
	}
	return s.store.DeletePassenger(ctx, id)
}
