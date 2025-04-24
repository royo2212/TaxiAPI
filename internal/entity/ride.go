package entity

type Ride struct {
	RideID      int        `json:"ride_id"`
	PassengerID int        `json:"-"`
	Passenger   *Passenger `json:"passenger,omitempty"`
	Driver      *Driver    `json:"driver,omitempty"`
	DriverID    int        `json:"-"`
	Origin      string     `json:"origin"`
	Destination string     `json:"destination"`
	Status      Status     `json:"status"`
}
type Status string

const (
	StatusPending   Status = "pending"
	StatusAccepted  Status = "accepted"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}
