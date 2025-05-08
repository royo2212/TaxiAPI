package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"taxiAPI/internal/entity"
	"taxiAPI/internal/service"
	customErrors "taxiAPI/internal/errors"
)

type DriverHandler struct {
	service *service.DriverService
}

func NewDriverHandler(service *service.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

type registerDriverRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  int    `json:"phone_number"`
	CarType      string `json:"car_type"`
	LicensePlate int    `json:"license_plate"`
}

func (h *DriverHandler) RegisterDriver(w http.ResponseWriter, r *http.Request) {
	var req registerDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	driver := &entity.Driver{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		CarType:      req.CarType,
		LicensePlate: req.LicensePlate,
		IsAvailable:  true,
	}

	created, err := h.service.RegisterDriver(r.Context(),driver)
	if err != nil {
		switch err {
        case customErrors.ErrDriverDataRequired:
            http.Error(w, err.Error(), http.StatusBadRequest)
        case customErrors.ErrPhoneNumberExists:
            http.Error(w, err.Error(), http.StatusConflict)
        case customErrors.ErrFirstName, customErrors.ErrLastName, 
             customErrors.ErrPhoneNumber, customErrors.ErrCarTypeRequired, 
             customErrors.ErrLicensePlateRequired:
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *DriverHandler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	driver, err := h.service.GetDriverByID(r.Context(),id)
	if err != nil {
		if err == customErrors.ErrDriverNotFound {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(driver)
	if err != nil {
		return
	}
}

func (h *DriverHandler) GetAllDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.service.GetAllDrivers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func (h *DriverHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteDriver(r.Context(), id)
	if err != nil {
		if err == customErrors.ErrDriverNotFound {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Driver deleted successfully",
		"id":      idStr,
	})
}
