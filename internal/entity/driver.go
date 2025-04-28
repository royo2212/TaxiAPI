package entity

type Driver struct {
	DriverID     int    `json:"driver_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  int    `json:"phone_number"`
	IsAvailable  bool   `json:"is_available"`
	CarType      string `json:"car_type"`
	LicensePlate int    `json:"license_plate"`
}
