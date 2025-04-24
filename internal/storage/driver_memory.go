package storage

import (
	"context"
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type DriverMemory struct {
	mutex   sync.RWMutex
	drivers map[int]*entity.Driver
	nextID  int
}

func NewDriverMemory() *DriverMemory {
	return &DriverMemory{
		drivers: make(map[int]*entity.Driver),
		nextID:  1,
	}
}

func (m *DriverMemory) RegisterDriver(ctx context.Context, d *entity.Driver) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		m.mutex.Lock()
		defer m.mutex.Unlock()
		d.DriverID = m.nextID
		m.drivers[d.DriverID] = d
		m.nextID++
		return d, nil
	}
}

func (m *DriverMemory) GetDriverByID(ctx context.Context, id int) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	p, ok := m.drivers[id]
	if !ok {
		return nil, customErrors.ErrDriverNotFound
	}
	return p, nil
}

func (m *DriverMemory) GetAllDrivers(ctx context.Context) ([]*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	drivers := make([]*entity.Driver, 0, len(m.drivers))
	for _, d := range m.drivers {
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (m *DriverMemory) DeleteDriver(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.drivers[id]
	if !ok {
		return customErrors.ErrDriverNotFound
	}
	delete(m.drivers, id)
	return nil
}
func (m *DriverMemory) FindByPhoneNumber(ctx context.Context, phone int) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, p := range m.drivers {
		if p.PhoneNumber == phone {
			return p, nil
		}
	}
	return nil, customErrors.ErrDriverNotFound
}
