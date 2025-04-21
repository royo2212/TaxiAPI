package entity

type Passenger struct {
	PassengerID int    `json:"passenger_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber int    `json:"phone_number"`
}
