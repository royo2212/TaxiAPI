package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"taxiAPI/internal/entity"
	"taxiAPI/internal/service"
)

type RideHandler struct {
	service *service.RideService
}

func NewRideHandler(service *service.RideService) *RideHandler {
	return &RideHandler{
		service: service,
	}
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type createRideRequest struct {
	PassengerID int    `json:"passenger_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func (h *RideHandler) CreateRide(w http.ResponseWriter, r *http.Request) {
	var req createRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	ride, err := h.service.CreateRide(req.PassengerID, req.Origin, req.Destination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(ride); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) GetRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	rideID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ride ID", http.StatusBadRequest)
		return
	}

	ride, err := h.service.GetRide(rideID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ride); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) GetAllRides(w http.ResponseWriter, r *http.Request) {
	rides, err := h.service.GetAllRides()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rides); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) AssignDriverToRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	rideID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ride ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DriverID int `json:"driver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.AssignDriverToRide(rideID, req.DriverID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Driver assigned successfully",
		"ride_id":   rideID,
		"driver_id": req.DriverID,
	})
	if err != nil {
		return
	}
}
func (h *RideHandler) UpdateRideStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	rideID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ride ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRideStatus(rideID, entity.Status(req.Status)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Ride status updated",
		"ride_id": rideID,
		"status":  req.Status,
	})
	if err != nil {
		return
	}
}
