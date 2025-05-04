package storage

import (
	"context"
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type Passenger struct {
	mutex      sync.RWMutex
	passengers map[int]*entity.Passenger
	nextID     int
}

func NewPassenger() *Passenger {
	return &Passenger{
		passengers: make(map[int]*entity.Passenger),
		nextID:     1,
	}
}

func (p *Passenger) RegisterPassenger(ctx context.Context, passenger *entity.Passenger) (*entity.Passenger, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p.mutex.Lock()
		defer p.mutex.Unlock()
		passenger.ID = p.nextID
		p.passengers[passenger.ID] = passenger
		p.nextID++
		return passenger, nil
	}
}

func (p *Passenger) GetPassengerByID(ctx context.Context, id int) (*entity.Passenger, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	passenger, ok := p.passengers[id]
	if !ok {
		return nil, customErrors.ErrPassengerNotFound
	}
	return passenger, nil
}

func (p *Passenger) GetAllPassengers(ctx context.Context) ([]*entity.Passenger, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	passengers := make([]*entity.Passenger, 0, len(p.passengers))
	for _, passenger := range p.passengers {
		passengers = append(passengers, passenger)
	}
	return passengers, nil
}

func (p *Passenger) DeletePassenger(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_, ok := p.passengers[id]
	if !ok {
		return customErrors.ErrPassengerNotFound
	}
	delete(p.passengers, id)
	return nil
}

func (p *Passenger) FindByPhoneNumber(ctx context.Context, phone int) (*entity.Passenger, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, passenger := range p.passengers {
		if passenger.PhoneNumber == phone {
			return passenger, nil
		}
	}
	return nil, customErrors.ErrPassengerNotFound
}
