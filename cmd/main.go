package main

import (
	"log"
	"net/http"
	"taxiAPI/config"
	"taxiAPI/internal/endpoints"
	"taxiAPI/internal/service"
	"taxiAPI/internal/storage/postgres"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()
	log.Println("Connected to Postgres")
	passengerStore := postgres.NewPassengerPostgres(db)
	driverStore := postgres.NewDriverPostgres(db)
	rideStore := postgres.NewRidePostgres(db)
	passengerService := service.NewPassengerService(passengerStore)
	driverService := service.NewDriverService(driverStore)
	rideService := service.NewRideService(rideStore, passengerStore, driverStore)
	passengerHandler := endpoints.NewPassengerHandler(passengerService)
	driverHandler := endpoints.NewDriverHandler(driverService)
	rideHandler := endpoints.NewRideHandler(rideService)
	router := mux.NewRouter()
	router.HandleFunc("/passengers", passengerHandler.RegisterPassenger).Methods("POST")
	router.HandleFunc("/passengers", passengerHandler.GetAllPassengers).Methods("GET")
	router.HandleFunc("/passengers/{id}", passengerHandler.GetPassengerByID).Methods("GET")
	router.HandleFunc("/passengers/{id}", passengerHandler.DeletePassenger).Methods("DELETE")
	router.HandleFunc("/drivers", driverHandler.RegisterDriver).Methods("POST")
	router.HandleFunc("/drivers", driverHandler.GetAllDrivers).Methods("GET")
	router.HandleFunc("/drivers/{id}", driverHandler.GetDriverByID).Methods("GET")
	router.HandleFunc("/drivers/{id}", driverHandler.DeleteDriver).Methods("DELETE")
	router.HandleFunc("/rides", rideHandler.CreateRide).Methods("POST")
	router.HandleFunc("/rides", rideHandler.GetAllRides).Methods("GET")
	router.HandleFunc("/rides/{id}", rideHandler.GetRide).Methods("GET")
	router.HandleFunc("/rides/{id}/driver", rideHandler.AssignDriverToRide).Methods("PUT")
	router.HandleFunc("/rides/{id}/status", rideHandler.UpdateRideStatus).Methods("PUT")
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
