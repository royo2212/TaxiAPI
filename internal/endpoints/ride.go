package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"taxiAPI/internal/entity"
	"taxiAPI/internal/service"
	customErrors "taxiAPI/internal/errors"
	"github.com/gorilla/mux"
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	ride, err := h.service.CreateRide(r.Context(), req.PassengerID, req.Origin, req.Destination)
	if err != nil {
		switch err {
			case customErrors.ErrPassengerIDRequired, customErrors.ErrOriginRequired, 
				 customErrors.ErrDestinationRequired:
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			case customErrors.ErrPassengerNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(ride); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) GetRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	rideID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ride ID", http.StatusBadRequest)
		return
	}

	ride, err := h.service.GetRide(r.Context(), rideID)
	if err != nil {
		if err == customErrors.ErrRideNotFound {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else if err == customErrors.ErrRideIDRequired {
            http.Error(w, err.Error(), http.StatusBadRequest)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ride); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) GetAllRides(w http.ResponseWriter, r *http.Request) {
	rides, err := h.service.GetAllRides(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rides); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RideHandler) AssignDriverToRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	rideID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ride ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DriverID int `json:"driver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.AssignDriverToRide(r.Context(), rideID, req.DriverID); err != nil {
		switch err {
        case customErrors.ErrRideNotFound, customErrors.ErrDriverNotFound:
            http.Error(w, err.Error(), http.StatusNotFound)
        case customErrors.ErrRideIDRequired, customErrors.ErrDriverIDRequired:
            http.Error(w, err.Error(), http.StatusBadRequest)
        case customErrors.ErrDriverAlreadyAssignedToRide, customErrors.ErrRideAlreadyAssigned:
            http.Error(w, err.Error(), http.StatusConflict)
        case customErrors.ErrCannotAssignDriverToNonPendingRide, customErrors.ErrDriverAlreadyOnActiveRide:
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
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
		http.Error(w, "Invalid ride ID", http.StatusBadRequest)
		return
	}
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRideStatus(r.Context(), rideID, entity.Status(req.Status)); err != nil {
		switch err {
        case customErrors.ErrRideNotFound:
            http.Error(w, err.Error(), http.StatusNotFound)
        case customErrors.ErrRideIDRequired:
            http.Error(w, err.Error(), http.StatusBadRequest)
        case customErrors.ErrInvalidRideStatus:
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        case customErrors.ErrCannotChangeCompletedRide:
            http.Error(w, err.Error(), http.StatusConflict)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
	w.Header().Set("Content-Type", "application/json")
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
