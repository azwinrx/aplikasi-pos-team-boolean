package adaptor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type StaffHandler struct {
	service usecase.StaffUseCase
}

func NewStaffAdaptor(service usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{service}
}

// GetList menangani GET /staff
func (h *StaffHandler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	role := r.URL.Query().Get("role")

	// Construct DTO
	req := dto.StaffFilterRequest{
		Page:      page,
		Limit:     limit,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Role:      role,
	}

	result, pagination, err := h.service.GetListStaff(ctx, req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", result, pagination)
}

// GetByID menangani GET /staff/{id}
func (h *StaffHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid staff id", nil)
		return
	}

	result, err := h.service.GetStaffByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get staff detail", result)
}

// Create menangani POST /staff
func (h *StaffHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req dto.StaffCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "validation error", err)
		return
	}

	result, err := h.service.CreateStaff(ctx, req)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.ResponseBadRequest(w, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success create staff", result)
}

// Update menangani PUT /staff/{id}
func (h *StaffHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid staff id", nil)
		return
	}

	// Parse request body
	var req dto.StaffUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "validation error", err)
		return
	}

	result, err := h.service.UpdateStaff(ctx, uint(id), req)
	if err != nil {
		if err.Error() == "staff not found" {
			utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "email already exists" {
			utils.ResponseBadRequest(w, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success update staff", result)
}

// Delete menangani DELETE /staff/{id}
func (h *StaffHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid staff id", nil)
		return
	}

	err = h.service.DeleteStaff(ctx, uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success delete staff", nil)
}
