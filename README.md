# 🚕 Mini Ride Booking API (Go)

A simple backend service for creating and managing taxi rides — built in Go using Clean Architecture.

---

## ✅ Features Implemented

### 🧍 Passenger
- ➕ Register a new passenger → `POST /passengers`
- 📋 Get all passengers → `GET /passengers`
- 🔍 Get passenger by ID → `GET /passengers/{id}`
- ❌ Delete passenger → `DELETE /passengers/{id}`

### 🚗 Driver
- ➕ Register a new driver → `POST /drivers`
- 📋 Get all drivers → `GET /drivers`
- 🔍 Get driver by ID → `GET /drivers/{id}`
- ❌ Delete driver → `DELETE /drivers/{id}`

### 🚕 Ride
- ➕ Create a new ride → `POST /rides`
- 📋 Get all rides → `GET /rides`
- 🔍 Get ride by ID → `GET /rides/{id}`
- 👨‍✈️ Assign a driver to a ride → `PUT /rides/{id}/driver`
- 🔄 Update ride status → `PUT /rides/{id}/status`

---

## ⚠️ Rules & Validations

- A **driver can only be assigned to one ride at a time** unless their current ride is `completed` or `cancelled`.
- Rides are automatically marked as `"accepted"` when a driver is assigned.
- Full `passenger` and `driver` data is returned inside each ride object.
- Phone numbers must be unique for both passengers and drivers.

---

## ▶️ How to Run

1. Clone the repo
2. Navigate to the project directory
3. Run the server:

```bash
go run ./cmd/main.go
```
---

## 📬 How to Use with Postman

Register a passenger with `POST /passengers`
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": 123456789
}
```

Register a driver with `POST /drivers`
```json
{
  "first_name": "Alex",
  "last_name": "Smith",
  "phone_number": 111222333,
  "car_type": "Toyota Prius",
  "license_plate": 123456
}
```

Create a ride with `POST /rides`
```json
{
  "passenger_id": 1,
  "origin": "Tel Aviv",
  "destination": "Jerusalem"
}
```

Assign a driver with `PUT /rides/1/driver`
```json
{
  "driver_id": 1
}
```

Update ride status with `PUT /rides/1/status`
```json
{
  "status": "completed"
}
```

Valid status values: `"pending"`, `"accepted"`, `"completed"`, `"cancelled"`  
Use `GET /rides/1` to fetch a specific ride (returns full passenger & driver).  
List everything with: `GET /rides`, `GET /passengers`, `GET /drivers`  
Delete with: `DELETE /passengers/{id}`, `DELETE /drivers/{id}`


