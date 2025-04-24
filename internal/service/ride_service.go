package service

import (
	"taxiAPI/internal/entity"
	customErrors "taxiAPI/internal/errors"
)

type RideStore interface {
	SaveRide(ride *entity.Ride) error
	FindRideByID(id int) (*entity.Ride, error)
	UpdateRideStatus(rideID int, status entity.Status) error
	AssignDriverToRide(rideID int, driverID int) error
	GetAllRides() ([]*entity.Ride, error)
	FindActiveRideByDriver(driverID int) (*entity.Ride, error)
}
type RideService struct {
	store          RideStore
	passengerStore PassengerStore
	driverStore    DriverStore
	nextID         int
}

func NewRideService(store RideStore, passengerStore PassengerStore, driverStore DriverStore) *RideService {
	return &RideService{
		store:          store,
		passengerStore: passengerStore,
		driverStore:    driverStore,
		nextID:         1,
	}
}
func (s *RideService) CreateRide(passengerID int, origin, destination string) (*entity.Ride, error) {
	if passengerID == 0 {
		return nil, customErrors.ErrPassengerIDRequired
	}
	if origin == "" {
		return nil, customErrors.ErrOriginRequired
	}
	if destination == "" {
		return nil, customErrors.ErrDestinationRequired
	}
	passenger, err := s.passengerStore.GetPassengerByID(passengerID)
	if err != nil {
		return nil, err
	}
	ride := &entity.Ride{
		RideID:      s.nextID,
		PassengerID: passengerID,
		Passenger:   passenger,
		Origin:      origin,
		Destination: destination,
		Status:      entity.StatusPending,
	}
	s.nextID++
	err = s.store.SaveRide(ride)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (s *RideService) GetRide(rideID int) (*entity.Ride, error) {
	if rideID == 0 {
		return nil, customErrors.ErrRideIDRequired
	}
	ride, err := s.store.FindRideByID(rideID)
	if err != nil {
		return nil, err
	}
	passenger, err := s.passengerStore.GetPassengerByID(ride.PassengerID)
	if err == nil {
		ride.Passenger = passenger
	}
	driver, err := s.driverStore.GetDriverByID(ride.DriverID)
	if err == nil {
		ride.Driver = driver
	}
	return ride, nil
}

func (s *RideService) GetAllRides() ([]*entity.Ride, error) {
	rides, err := s.store.GetAllRides()
	if err != nil {
		return nil, err
	}

	for _, ride := range rides {
		if passenger, err := s.passengerStore.GetPassengerByID(ride.PassengerID); err == nil {
			ride.Passenger = passenger
		}
		if driver, err := s.driverStore.GetDriverByID(ride.DriverID); err == nil {
			ride.Driver = driver
		}
	}

	return rides, nil
}

func (s *RideService) UpdateRideStatus(rideID int, status entity.Status) error {
	if rideID == 0 {
		return customErrors.ErrRideIDRequired
	}
	ride, err := s.store.FindRideByID(rideID)
	if err != nil {
		return err
	}
	if ride.Status == entity.StatusCompleted && status != entity.StatusCompleted {
		return customErrors.ErrCannotChangeCompletedRide
	}
	if !status.IsValid() {
		return customErrors.ErrInvalidRideStatus
	}
	ride.Status = status
	return s.store.UpdateRideStatus(rideID, ride.Status)
}

func (s *RideService) AssignDriverToRide(rideID, driverID int) error {
	if rideID == 0 {
		return customErrors.ErrRideIDRequired
	}
	if driverID == 0 {
		return customErrors.ErrDriverIDRequired
	}
	ride, err := s.store.FindRideByID(rideID)
	if err != nil {
		return err
	}
	if ride.DriverID != 0 {
		if ride.DriverID == driverID {
			return customErrors.ErrDriverAlreadyAssignedToRide
		}
		return customErrors.ErrRideAlreadyAssigned
	}
	if ride.Status != entity.StatusPending {
		return customErrors.ErrCannotAssignDriverToNonPendingRide
	}
	existingRide, err := s.store.FindActiveRideByDriver(driverID)
	if err == nil && existingRide != nil {
		return customErrors.ErrDriverAlreadyOnActiveRide
	}
	ride.DriverID = driverID
	ride.Status = entity.StatusAccepted
	return s.store.AssignDriverToRide(rideID, driverID)
}
