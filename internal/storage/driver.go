package storage

import (
	"context"
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type Driver struct {
	mutex   sync.RWMutex
	drivers map[int]*entity.Driver
	nextID  int
}

func NewDriver() *Driver {
	return &Driver{
		drivers: make(map[int]*entity.Driver),
		nextID:  1,
	}
}

func (d *Driver) RegisterDriver(ctx context.Context, driver *entity.Driver) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		d.mutex.Lock()
		defer d.mutex.Unlock()
		driver.DriverID = d.nextID
		d.drivers[driver.DriverID] = driver
		d.nextID++
		return driver, nil
	}
}

func (d *Driver) GetDriverByID(ctx context.Context, id int) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	driver, ok := d.drivers[id]
	if !ok {
		return nil, customErrors.ErrDriverNotFound
	}
	return driver, nil
}

func (d *Driver) GetAllDrivers(ctx context.Context) ([]*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	drivers := make([]*entity.Driver, 0, len(d.drivers))
	for _, drv := range d.drivers {
		drivers = append(drivers, drv)
	}
	return drivers, nil
}

func (d *Driver) DeleteDriver(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.drivers[id]
	if !ok {
		return customErrors.ErrDriverNotFound
	}
	delete(d.drivers, id)
	return nil
}

func (d *Driver) FindByPhoneNumber(ctx context.Context, phone int) (*entity.Driver, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, driver := range d.drivers {
		if driver.PhoneNumber == phone {
			return driver, nil
		}
	}
	return nil, customErrors.ErrDriverNotFound
}
