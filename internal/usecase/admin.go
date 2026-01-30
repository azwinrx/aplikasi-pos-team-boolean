package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/pkg/utils"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// AdminUseCase mendefinisikan interface untuk admin management
type AdminUseCase interface {
	ListAdmins(ctx context.Context, page int, limit int, role string) (*dto.ListAdminResponse, error)
	EditAdminAccess(ctx context.Context, adminID uint, req *dto.EditAdminAccessRequest) (*dto.EditAdminAccessResponse, error)
	CreateAdminWithEmail(ctx context.Context, req *dto.CreateAdminRequest) (*dto.AdminResponse, error)
	UpdateUserProfile(ctx context.Context, userID uint, req *dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error)
	GetUserProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error)
	Logout(ctx context.Context, userID uint) error
}

// adminUseCase implementasi dari AdminUseCase interface
type adminUseCase struct {
	authRepo      repository.AuthRepository
	emailService  *utils.EmailService
	logger        *zap.Logger
}

// NewAdminUseCase membuat instance baru dari adminUseCase
func NewAdminUseCase(authRepo repository.AuthRepository, emailService *utils.EmailService, logger *zap.Logger) AdminUseCase {
	return &adminUseCase{
		authRepo:     authRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// ListAdmins mengambil daftar admin dengan pagination
func (u *adminUseCase) ListAdmins(ctx context.Context, page int, limit int, role string) (*dto.ListAdminResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Get admins from repository
	admins, total, err := u.authRepo.GetAdminsList(ctx, offset, limit, role)
	if err != nil {
		u.logger.Error("Failed to get admins list",
			zap.Error(err),
		)
		return nil, err
	}

	// Convert to response DTO
	adminResponses := make([]dto.AdminResponse, len(admins))
	for i, admin := range admins {
		adminResponses[i] = dto.AdminResponse{
			ID:        admin.ID,
			Email:     admin.Email,
			Name:      admin.Name,
			Role:      admin.Role,
			Status:    admin.Status,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		}
	}

	totalPages := (int(total) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	response := &dto.ListAdminResponse{
		Data:       adminResponses,
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	u.logger.Info("Admins list retrieved successfully",
		zap.Int("count", len(admins)),
		zap.Int64("total", total),
	)

	return response, nil
}

// EditAdminAccess mengedit akses admin (hanya untuk superadmin)
func (u *adminUseCase) EditAdminAccess(ctx context.Context, adminID uint, req *dto.EditAdminAccessRequest) (*dto.EditAdminAccessResponse, error) {
	// Validasi input
	if adminID == 0 {
		u.logger.Warn("Invalid admin ID",
			zap.Uint("admin_id", adminID),
		)
		return nil, errors.New("admin_id tidak valid")
	}

	// Validasi role
	validRoles := map[string]bool{"admin": true, "superadmin": true, "user": true}
	if !validRoles[req.Role] {
		u.logger.Warn("Invalid role",
			zap.String("role", req.Role),
		)
		return nil, fmt.Errorf("role '%s' tidak valid. gunakan 'admin', 'superadmin', atau 'user'", req.Role)
	}

	// Validasi status
	validStatus := map[string]bool{"active": true, "inactive": true}
	if !validStatus[req.Status] {
		u.logger.Warn("Invalid status",
			zap.String("status", req.Status),
		)
		return nil, fmt.Errorf("status '%s' tidak valid. gunakan 'active' atau 'inactive'", req.Status)
	}

	// Get admin untuk verifikasi
	admin, err := u.authRepo.GetUserByID(ctx, adminID)
	if err != nil {
		u.logger.Error("Failed to get admin",
			zap.Uint("admin_id", adminID),
			zap.Error(err),
		)
		return nil, err
	}

	if admin == nil {
		u.logger.Warn("Admin not found",
			zap.Uint("admin_id", adminID),
		)
		return nil, errors.New("admin tidak ditemukan")
	}

	// Prevent deactivating the only superadmin
	if req.Status == "inactive" && admin.Role == "superadmin" {
		count, err := u.authRepo.CountSuperadmins(ctx)
		if err != nil {
			u.logger.Error("Failed to count superadmins",
				zap.Error(err),
			)
			return nil, err
		}

		if count <= 1 {
			u.logger.Warn("Cannot deactivate the only superadmin")
			return nil, errors.New("tidak dapat menonaktifkan satu-satunya superadmin")
		}
	}

	// Update admin
	admin.Role = req.Role
	admin.Status = req.Status

	if err := u.authRepo.UpdateUser(ctx, admin); err != nil {
		u.logger.Error("Failed to update admin",
			zap.Uint("admin_id", adminID),
			zap.Error(err),
		)
		return nil, err
	}

	response := &dto.EditAdminAccessResponse{
		ID:      admin.ID,
		Email:   admin.Email,
		Name:    admin.Name,
		Role:    admin.Role,
		Status:  admin.Status,
		Message: fmt.Sprintf("Akses admin %s berhasil diubah", admin.Email),
	}

	u.logger.Info("Admin access updated successfully",
		zap.Uint("admin_id", adminID),
		zap.String("new_role", req.Role),
		zap.String("new_status", req.Status),
	)

	return response, nil
}

// CreateAdminWithEmail membuat admin baru dan mengirim password via email
func (u *adminUseCase) CreateAdminWithEmail(ctx context.Context, req *dto.CreateAdminRequest) (*dto.AdminResponse, error) {
	// Validasi input
	if req.Email == "" || req.Name == "" {
		u.logger.Warn("Missing required fields",
			zap.String("email", req.Email),
			zap.String("name", req.Name),
		)
		return nil, errors.New("email dan name harus diisi")
	}

	// Validasi role
	validRoles := map[string]bool{"admin": true, "superadmin": true}
	if !validRoles[req.Role] {
		u.logger.Warn("Invalid role",
			zap.String("role", req.Role),
		)
		return nil, fmt.Errorf("role hanya bisa 'admin' atau 'superadmin'")
	}

	// Check if email already exists
	existing, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Failed to check existing email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	if existing != nil {
		u.logger.Warn("Email already exists",
			zap.String("email", req.Email),
		)
		return nil, errors.New("email sudah terdaftar")
	}

	// Generate random password
	generatedPassword := u.generateRandomPassword(12)

	// Hash password
	hashedPassword, err := utils.HashPassword(generatedPassword)
	if err != nil {
		u.logger.Error("Failed to hash password",
			zap.Error(err),
		)
		return nil, err
	}

	// Create user
	user := &entity.User{
		Email:    strings.ToLower(req.Email),
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
		Status:   "active",
	}

	if err := u.authRepo.CreateUser(ctx, user); err != nil {
		u.logger.Error("Failed to create admin user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	// Send email with generated password
	emailBody := fmt.Sprintf(`
Halo %s,

Akun admin Anda telah berhasil dibuat di sistem POS.

Berikut adalah kredensial akun Anda:
Email: %s
Password: %s
Role: %s

Silakan login dan ubah password Anda di halaman profil.
Jangan bagikan password ini kepada orang lain.

Best regards,
POS System Administrator
	`, user.Name, user.Email, generatedPassword, user.Role)

	if err := u.emailService.SendEmail(user.Email, "Admin Account Created", emailBody); err != nil {
		u.logger.Error("Failed to send password email",
			zap.String("email", user.Email),
			zap.Error(err),
		)
		// Log error but don't fail - user is already created
		// In production, you might want to handle this differently
	}

	response := &dto.AdminResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	u.logger.Info("Admin created successfully and password sent via email",
		zap.Uint("admin_id", user.ID),
		zap.String("email", user.Email),
	)

	return response, nil
}

// UpdateUserProfile mengupdate profil user
func (u *adminUseCase) UpdateUserProfile(ctx context.Context, userID uint, req *dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error) {
	// Validasi input
	if userID == 0 {
		u.logger.Warn("Invalid user ID",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user_id tidak valid")
	}

	// Get user
	user, err := u.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to get user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	if user == nil {
		u.logger.Warn("User not found",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user tidak ditemukan")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		// Hash new password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			u.logger.Error("Failed to hash password",
				zap.Error(err),
			)
			return nil, err
		}
		user.Password = hashedPassword
	}

	// Update in database
	if err := u.authRepo.UpdateUser(ctx, user); err != nil {
		u.logger.Error("Failed to update user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	response := &dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Message:   "Profil berhasil diubah",
	}

	u.logger.Info("User profile updated successfully",
		zap.Uint("user_id", userID),
	)

	return response, nil
}

// GetUserProfile mengambil profil user
func (u *adminUseCase) GetUserProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error) {
	if userID == 0 {
		u.logger.Warn("Invalid user ID",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user_id tidak valid")
	}

	user, err := u.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to get user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	if user == nil {
		u.logger.Warn("User not found",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user tidak ditemukan")
	}

	response := &dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	u.logger.Info("User profile retrieved successfully",
		zap.Uint("user_id", userID),
	)

	return response, nil
}

// Logout endpoint (just log the action)
func (u *adminUseCase) Logout(ctx context.Context, userID uint) error {
	u.logger.Info("User logged out successfully",
		zap.Uint("user_id", userID),
	)

	// In a real app, you might want to invalidate tokens, update last_login, etc.
	// For now, just log the action
	return nil
}

// generateRandomPassword generates a random password
func (u *adminUseCase) generateRandomPassword(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
