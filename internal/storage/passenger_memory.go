package storage

import (
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type PassengerMemory struct {
	mutex      sync.RWMutex
	passengers map[int]*entity.Passenger
	nextID     int
}

func NewPassengerMemory() *PassengerMemory {
	return &PassengerMemory{
		passengers: make(map[int]*entity.Passenger),
		nextID:     1,
	}
}

func (m *PassengerMemory) RegisterPassenger(p *entity.Passenger) (*entity.Passenger, error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()
	p.PassengerID = m.nextID
	m.passengers[p.PassengerID] = p
	m.nextID++

	return p, nil
}

func (m *PassengerMemory) GetPassengerByID(id int) (*entity.Passenger, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	p, ok := m.passengers[id]
	if !ok {
		return nil, customErrors.ErrPassengerNotFound
	}
	return p, nil
}

func (m *PassengerMemory) GetAllPassengers() ([]*entity.Passenger, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	passengers := make([]*entity.Passenger, 0, len(m.passengers))
	for _, p := range m.passengers {
		passengers = append(passengers, p)
	}
	return passengers, nil
}

func (m *PassengerMemory) DeletePassenger(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.passengers[id]
	if !ok {
		return customErrors.ErrPassengerNotFound
	}
	delete(m.passengers, id)
	return nil
}
func (m *PassengerMemory) FindByPhoneNumber(phone int) (*entity.Passenger, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, p := range m.passengers {
		if p.PhoneNumber == phone {
			return p, nil
		}
	}
	return nil, customErrors.ErrPassengerNotFound
}
