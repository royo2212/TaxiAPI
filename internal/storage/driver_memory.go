package storage

import (
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

func (m *DriverMemory) RegisterDriver(d *entity.Driver) (*entity.Driver, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	d.DriverID = m.nextID
	m.drivers[d.DriverID] = d
	m.nextID++
	return d, nil
}

func (m *DriverMemory) GetDriverByID(id int) (*entity.Driver, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	p, ok := m.drivers[id]
	if !ok {
		return nil, customErrors.ErrDriverNotFound
	}
	return p, nil
}

func (m *DriverMemory) GetAllDrivers() ([]*entity.Driver, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	drivers := make([]*entity.Driver, 0, len(m.drivers))
	for _, d := range m.drivers {
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (m *DriverMemory) DeleteDriver(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.drivers[id]
	if !ok {
		return customErrors.ErrDriverNotFound
	}
	delete(m.drivers, id)
	return nil
}
func (m *DriverMemory) FindByPhoneNumber(phone int) (*entity.Driver, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, p := range m.drivers {
		if p.PhoneNumber == phone {
			return p, nil
		}
	}
	return nil, customErrors.ErrDriverNotFound
}
