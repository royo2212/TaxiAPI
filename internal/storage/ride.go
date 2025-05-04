package storage

import (
	"context"
	"sync"
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type Ride struct {
	mutex  sync.RWMutex
	rides  map[int]*entity.Ride
	nextID int
}

func NewRide() *Ride {
	return &Ride{
		rides:  make(map[int]*entity.Ride),
		nextID: 1,
	}
}

func (r *Ride) SaveRide(ctx context.Context, ride *entity.Ride) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		ride.RideID = r.nextID
		r.rides[ride.RideID] = ride
		r.nextID++
		return nil
	}
}

func (r *Ride) FindRideByID(ctx context.Context, id int) (*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ride, ok := r.rides[id]
	if !ok {
		return nil, customErrors.ErrRideNotFound
	}
	return ride, nil
}

func (r *Ride) UpdateRideStatus(ctx context.Context, rideID int, status entity.Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ride, ok := r.rides[rideID]
	if !ok {
		return customErrors.ErrRideNotFound
	}
	ride.Status = status
	return nil
}

func (r *Ride) AssignDriverToRide(ctx context.Context, rideID, driverID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	ride, ok := r.rides[rideID]
	if !ok {
		return customErrors.ErrRideNotFound
	}

	ride.DriverID = &driverID
	return nil
}

func (r *Ride) GetAllRides(ctx context.Context) ([]*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	rides := make([]*entity.Ride, 0, len(r.rides))
	for _, ride := range r.rides {
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *Ride) FindActiveRideByDriver(ctx context.Context, driverID int) (*entity.Ride, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, ride := range r.rides {
		if ride.DriverID != nil &&
			*ride.DriverID == driverID &&
			ride.Status != entity.StatusCompleted &&
			ride.Status != entity.StatusCancelled {
			return ride, nil
		}
	}
	return nil, customErrors.ErrRideNotFound
}
