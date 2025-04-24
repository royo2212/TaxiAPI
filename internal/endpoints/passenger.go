package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"taxiAPI/internal/entity"
	"taxiAPI/internal/service"
)

type PassengerHandler struct {
	service *service.PassengerService
}

func NewPassengerHandler(service *service.PassengerService) *PassengerHandler {
	return &PassengerHandler{
		service: service,
	}
}

type registerPassengerRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber int    `json:"phone_number"`
}

func (h *PassengerHandler) RegisterPassenger(w http.ResponseWriter, r *http.Request) {
	var req registerPassengerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	passenger := &entity.Passenger{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}
	created, err := h.service.RegisterPassenger(passenger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
func (h *PassengerHandler) GetPassengerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid passenger ID", http.StatusBadRequest)
		return
	}
	passenger, err := h.service.GetPassengerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(passenger); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
func (h *PassengerHandler) GetAllPassengers(w http.ResponseWriter, r *http.Request) {
	passengers, err := h.service.GetAllPassengers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(passengers); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
func (h *PassengerHandler) DeletePassenger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid passenger ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeletePassenger(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
