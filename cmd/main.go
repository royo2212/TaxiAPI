package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"taxiAPI/internal/endpoints"
	"taxiAPI/internal/service"
	"taxiAPI/internal/storage"
)

func main() {
	// ✅ Initialize in-memory storage
	rideStore := storage.NewRideMemory()
	passengerStore := storage.NewPassengerMemory()
	driverStore := storage.NewDriverMemory()

	// ✅ Initialize services
	passengerService := service.NewPassengerService(passengerStore)
	rideService := service.NewRideService(rideStore, passengerStore, driverStore)
	driverService := service.NewDriverService(driverStore)

	// ✅ Initialize handlers
	passengerHandler := endpoints.NewPassengerHandler(passengerService)
	rideHandler := endpoints.NewRideHandler(rideService)
	driverHandler := endpoints.NewDriverHandler(driverService)

	// ✅ Setup router
	router := mux.NewRouter()

	// 🧍 Passenger routes
	router.HandleFunc("/passengers", passengerHandler.RegisterPassenger).Methods("POST")
	router.HandleFunc("/passengers", passengerHandler.GetAllPassengers).Methods("GET")
	router.HandleFunc("/passengers/{id}", passengerHandler.GetPassengerByID).Methods("GET")
	router.HandleFunc("/passengers/{id}", passengerHandler.DeletePassenger).Methods("DELETE")
	// 🚗 Driver routes
	router.HandleFunc("/drivers", driverHandler.RegisterDriver).Methods("POST")
	router.HandleFunc("/drivers", driverHandler.GetAllDrivers).Methods("GET")
	router.HandleFunc("/drivers/{id}", driverHandler.GetDriverByID).Methods("GET")
	router.HandleFunc("/drivers/{id}", driverHandler.DeleteDriver).Methods("DELETE")
	// 🚕 Ride routes
	router.HandleFunc("/rides", rideHandler.CreateRide).Methods("POST")
	router.HandleFunc("/rides", rideHandler.GetAllRides).Methods("GET")
	router.HandleFunc("/rides/{id}", rideHandler.GetRide).Methods("GET")
	router.HandleFunc("/rides/{id}/driver", rideHandler.AssignDriverToRide).Methods("PUT")
	router.HandleFunc("/rides/{id}/status", rideHandler.UpdateRideStatus).Methods("PUT")
	// ✅ Start server
	log.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
