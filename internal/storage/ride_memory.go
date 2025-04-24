package storage

import (
	"context"
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type RideMemory struct {
	mutex  sync.RWMutex
	rides  map[int]*entity.Ride
	nextID int
}

func NewRideMemory() *RideMemory {
	return &RideMemory{
		rides:  make(map[int]*entity.Ride),
		nextID: 1,
	}
}

func (store *RideMemory) SaveRide(ctx context.Context, ride *entity.Ride) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		store.mutex.Lock()
		defer store.mutex.Unlock()
		ride.RideID = store.nextID
		store.rides[ride.RideID] = ride
		store.nextID++
		return nil
	}
}
func (store *RideMemory) FindRideByID(ctx context.Context, id int) (*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
		store.mutex.RLock()
		defer store.mutex.RUnlock()
		ride, ok := store.rides[id]
		if !ok {
			return nil, customErrors.ErrRideNotFound
		}
		return ride, nil
}
func (store *RideMemory) UpdateRideStatus(ctx context.Context,rideID int, status entity.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	ride, ok := store.rides[rideID]
	if !ok {
		return customErrors.ErrRideNotFound
	}
	ride.Status = status
	return nil
}
func (store *RideMemory) AssignDriverToRide(ctx context.Context,rideID, driverID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	ride, ok := store.rides[rideID]
	if !ok {
		return customErrors.ErrRideNotFound
	}
	ride.DriverID = driverID
	return nil
}
func (store *RideMemory) GetAllRides(ctx context.Context) ([]*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	rides := make([]*entity.Ride, 0, len(store.rides))
	for _, ride := range store.rides {
		rides = append(rides, ride)
	}
	return rides, nil
}
func (store *RideMemory) FindActiveRideByDriver(ctx context.Context,driverID int) (*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, ride := range store.rides {
		if ride.DriverID == driverID && ride.Status != entity.StatusCompleted && ride.Status != entity.StatusCancelled {
			return ride, nil
		}
	}
	return nil, customErrors.ErrRideNotFound
}
