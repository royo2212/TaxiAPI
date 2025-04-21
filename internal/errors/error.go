package errors

import "errors"

var (
	ErrRideNotFound                       = errors.New("ride not found")
	ErrRideIDRequired                     = errors.New("ride ID is required")
	ErrOriginRequired                     = errors.New("origin is required")
	ErrDestinationRequired                = errors.New("destination is required")
	ErrPassengerIDRequired                = errors.New("passenger ID is required")
	ErrInvalidRideStatus                  = errors.New("invalid ride status")
	ErrCannotChangeCompletedRide          = errors.New("cannot change completed ride")
	ErrDriverIDRequired                   = errors.New("drive ID is required")
	ErrCannotAssignDriverToNonPendingRide = errors.New("cannot assign driver to non pending ride")
	ErrRideAlreadyAssigned                = errors.New("ride already assigned")
	ErrDriverAlreadyAssignedToRide        = errors.New("this driver already assigned to this ride")
	ErrPassengerNotFound                  = errors.New("passenger not found")
	ErrFirstName                          = errors.New("first name is required")
	ErrLastName                           = errors.New("last name is required")
	ErrPhoneNumber                        = errors.New("phone number is required")
	ErrPassengerDataRequired              = errors.New("passenger data required")
	ErrPhoneNumberExists                  = errors.New("phone number exists")
	ErrDriverNotFound                     = errors.New("driver not found")
	ErrDriverDataRequired                 = errors.New("driver data is required")
	ErrCarTypeRequired                    = errors.New("car type is required")
	ErrLicensePlateRequired               = errors.New("license plate is required")
	ErrDriverAlreadyOnActiveRide          = errors.New("driver already on active ride")
)
