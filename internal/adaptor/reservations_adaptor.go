package adaptor

import (
	"net/http"
	"strconv"

	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReservationsAdaptor menangani semua request HTTP untuk reservations
type ReservationsAdaptor struct {
	reservationsUsecase usecase.ReservationsUseCase
	logger              *zap.Logger
}

// NewReservationsAdaptor membuat instance baru dari ReservationsAdaptor
func NewReservationsAdaptor(reservationsUsecase usecase.ReservationsUseCase, logger *zap.Logger) *ReservationsAdaptor {
	return &ReservationsAdaptor{
		reservationsUsecase: reservationsUsecase,
		logger:              logger,
	}
}

// GetAllReservations menangani request untuk mengambil semua reservasi
func (h *ReservationsAdaptor) GetAllReservations(c *gin.Context) {
	h.logger.Debug("GetAllReservations handler called")

	response, err := h.reservationsUsecase.GetAllReservations(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get all reservations",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data reservasi: "+err.Error())
		return
	}

	h.logger.Info("GetAllReservations completed successfully",
		zap.Int("total_reservations", len(response)),
	)

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data reservasi berhasil diambil", response)
}

// CreateReservation menangani request untuk membuat reservasi baru
func (h *ReservationsAdaptor) CreateReservation(c *gin.Context) {
	h.logger.Debug("CreateReservation handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.ReservationCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create reservation",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.reservationsUsecase.CreateReservation(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create reservation",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal membuat reservasi: "+err.Error())
		return
	}

	h.logger.Info("Reservation created successfully",
		zap.Int64("id", response.ID),
	)

	utils.ResponseSuccess(c.Writer, http.StatusCreated, "Reservasi berhasil dibuat", response)
}

// UpdateReservation menangani request untuk update reservasi yang sudah ada
func (h *ReservationsAdaptor) UpdateReservation(c *gin.Context) {
	h.logger.Debug("UpdateReservation handler called", zap.String("client_ip", c.ClientIP()))

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	var req dto.ReservationUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update reservation",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	err = h.reservationsUsecase.UpdateReservation(c.Request.Context(), uint(id), req)
	if err != nil {
		h.logger.Error("Failed to update reservation",
			zap.Error(err),
			zap.Uint("id", uint(id)),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal memperbarui reservasi: "+err.Error())
		return
	}

	h.logger.Info("Reservation updated successfully", zap.Uint("id", uint(id)))

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Reservasi berhasil diperbarui", nil)
}

// DeleteReservation menangani request untuk menghapus reservasi
func (h *ReservationsAdaptor) DeleteReservation(c *gin.Context) {
	h.logger.Debug("DeleteReservation handler called", zap.String("client_ip", c.ClientIP()))

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	err = h.reservationsUsecase.DeleteReservation(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete reservation",
			zap.Error(err),
			zap.Uint("id", uint(id)),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal menghapus reservasi: "+err.Error())
		return
	}

	h.logger.Info("Reservation deleted successfully", zap.Uint("id", uint(id)))

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Reservasi berhasil dihapus", nil)
}
