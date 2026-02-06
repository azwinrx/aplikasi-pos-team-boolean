package adaptor

import (
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminAdaptor menghandle admin-related HTTP requests
type AdminAdaptor struct {
	adminUseCase usecase.AdminUseCase
	logger       *zap.Logger
}

// NewAdminAdaptor membuat instance baru dari AdminAdaptor
func NewAdminAdaptor(adminUseCase usecase.AdminUseCase, logger *zap.Logger) *AdminAdaptor {
	return &AdminAdaptor{
		adminUseCase: adminUseCase,
		logger:       logger,
	}
}

// GetUserProfile mengambil profil user yang login
// GET /api/v1/profile
func (a *AdminAdaptor) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("user_id tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := userID.(uint)
	if !ok {
		a.logger.Warn("user_id type assertion failed")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response, err := a.adminUseCase.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		a.logger.Error("Failed to get user profile",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Profil user berhasil diambil", response)
}

// UpdateUserProfile mengupdate profil user yang login
// PUT /api/v1/profile
func (a *AdminAdaptor) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("user_id tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := userID.(uint)
	if !ok {
		a.logger.Warn("user_id type assertion failed")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("Invalid request body",
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	response, err := a.adminUseCase.UpdateUserProfile(c.Request.Context(), id, &req)
	if err != nil {
		a.logger.Error("Failed to update user profile",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Profil berhasil diubah", response)
}

// ListAdmins mengambil daftar admin (hanya superadmin)
// GET /api/v1/admin/list
func (a *AdminAdaptor) ListAdmins(c *gin.Context) {
	// Check if user is superadmin
	userRole, exists := c.Get("user_role")
	if !exists {
		a.logger.Warn("user_role tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	role, ok := userRole.(string)
	if !ok || role != "superadmin" {
		a.logger.Warn("User is not superadmin",
			zap.String("role", role),
		)
		utils.ResponseError(c.Writer, http.StatusForbidden, "Anda tidak memiliki akses untuk melihat daftar admin")
		return
	}

	// Get query parameters
	page := 1
	limit := 10
	roleFilter := ""

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if role := c.Query("role"); role != "" {
		roleFilter = role
	}

	response, err := a.adminUseCase.ListAdmins(c.Request.Context(), page, limit, roleFilter)
	if err != nil {
		a.logger.Error("Failed to list admins",
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Daftar admin berhasil diambil", response)
}

// EditAdminAccess mengedit akses admin (hanya superadmin)
// PUT /api/v1/admin/:id/access
func (a *AdminAdaptor) EditAdminAccess(c *gin.Context) {
	// Check if user is superadmin
	userRole, exists := c.Get("user_role")
	if !exists {
		a.logger.Warn("user_role tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	role, ok := userRole.(string)
	if !ok || role != "superadmin" {
		a.logger.Warn("User is not superadmin",
			zap.String("role", role),
		)
		utils.ResponseError(c.Writer, http.StatusForbidden, "Hanya superadmin yang dapat mengedit akses admin")
		return
	}

	// Get admin ID from URL
	adminIDStr := c.Param("id")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 32)
	if err != nil {
		a.logger.Warn("Invalid admin ID",
			zap.String("admin_id", adminIDStr),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Admin ID tidak valid")
		return
	}

	// Parse request body
	var req dto.EditAdminAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("Invalid request body",
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	// Call usecase
	response, err := a.adminUseCase.EditAdminAccess(c.Request.Context(), uint(adminID), &req)
	if err != nil {
		a.logger.Error("Failed to edit admin access",
			zap.Uint("admin_id", uint(adminID)),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Akses admin berhasil diubah", response)
}

// CreateAdmin membuat admin baru (hanya superadmin)
// POST /api/v1/admin/create
func (a *AdminAdaptor) CreateAdmin(c *gin.Context) {
	// Check if user is superadmin
	userRole, exists := c.Get("user_role")
	if !exists {
		a.logger.Warn("user_role tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	role, ok := userRole.(string)
	if !ok || role != "superadmin" {
		a.logger.Warn("User is not superadmin",
			zap.String("role", role),
		)
		utils.ResponseError(c.Writer, http.StatusForbidden, "Hanya superadmin yang dapat membuat admin baru")
		return
	}

	// Parse request body
	var req dto.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("Invalid request body",
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	response, err := a.adminUseCase.CreateAdminWithEmail(c.Request.Context(), &req)
	if err != nil {
		a.logger.Error("Failed to create admin",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusCreated, "Admin berhasil dibuat. Password telah dikirim ke email", response)
}

// CreateAdminWithEmail is an alias for CreateAdmin
// POST /api/v1/admin
func (a *AdminAdaptor) CreateAdminWithEmail(c *gin.Context) {
	a.CreateAdmin(c)
}

// Logout melakukan logout
// POST /api/v1/auth/logout
func (a *AdminAdaptor) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("user_id tidak ditemukan di context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := userID.(uint)
	if !ok {
		a.logger.Warn("user_id type assertion failed")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := a.adminUseCase.Logout(c.Request.Context(), id)
	if err != nil {
		a.logger.Error("Failed to logout",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Logout berhasil", nil)
}
