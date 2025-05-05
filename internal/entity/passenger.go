package entity

type Passenger struct {
	ID int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber int    `json:"phone_number"`
}
