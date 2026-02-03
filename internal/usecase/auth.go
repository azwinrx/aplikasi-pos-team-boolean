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
	"time"

	"go.uber.org/zap"
)

// AuthUseCase mendefinisikan interface untuk auth business logic
type AuthUseCase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	GetUserByID(ctx context.Context, userID uint) (*dto.GetUserResponse, error)
	DeleteUser(ctx context.Context, userID uint) (*dto.DeleteUserResponse, error)
	CheckEmail(ctx context.Context, req dto.CheckEmailRequest) (*dto.CheckEmailResponse, error)
	SendOTP(ctx context.Context, req dto.SendOTPRequest) (*dto.SendOTPResponse, error)
	ValidateOTP(ctx context.Context, req dto.ValidateOTPRequest) (*dto.ValidateOTPResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

// authUsecase implementasi dari AuthUseCase interface
type authUsecase struct {
	authRepo     repository.AuthRepository
	logger       *zap.Logger
	emailService *utils.EmailService
}

// NewAuthUseCase membuat instance baru dari authUsecase
func NewAuthUseCase(authRepo repository.AuthRepository, logger *zap.Logger, emailService *utils.EmailService) AuthUseCase {
	return &authUsecase{
		authRepo:     authRepo,
		logger:       logger,
		emailService: emailService,
	}
}

// Login memproses login user
func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u.logger.Debug("Login attempt", zap.String("email", req.Email))

	// Get user dari database
	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Database error during login",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if user == nil {
		u.logger.Warn("Login failed - user not found",
			zap.String("email", req.Email),
		)
		return nil, errors.New("invalid email or password")
	}

	// Check if user is deleted
	if user.IsDeleted {
		u.logger.Warn("Login failed - user is deleted",
			zap.String("email", req.Email),
			zap.Uint("id", user.ID),
		)
		return nil, errors.New("your account has been deactivated")
	}

	// Validasi password
	if !utils.VerifyPassword(user.Password, req.Password) {
		u.logger.Warn("Login failed - invalid password",
			zap.String("email", req.Email),
		)
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		u.logger.Error("Failed to generate token",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to generate token")
	}

	u.logger.Info("Login successful",
		zap.String("email", req.Email),
		zap.Uint("user_id", user.ID),
	)

	return &dto.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// Register membuat user baru (registrasi)
func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	u.logger.Debug("Register attempt", zap.String("email", req.Email))

	// Check if email already exists
	existingUser, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Database error during registration",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if existingUser != nil && !existingUser.IsDeleted {
		u.logger.Warn("Registration failed - email already exists",
			zap.String("email", req.Email),
		)
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		u.logger.Error("Failed to hash password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to create user")
	}

	// Create user entity
	user := &entity.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Name:      req.Name,
		Role:      "customer", // Always set to customer
		Status:    "active",
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := u.authRepo.CreateUser(ctx, user); err != nil {
		u.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to create user")
	}

	u.logger.Info("User registered successfully",
		zap.String("email", req.Email),
		zap.Uint("user_id", user.ID),
	)

	return &dto.RegisterResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Role:    user.Role,
		Message: "User registered successfully",
	}, nil
}

// GetUserByID mendapatkan user berdasarkan ID
func (u *authUsecase) GetUserByID(ctx context.Context, userID uint) (*dto.GetUserResponse, error) {
	u.logger.Debug("Get user by ID", zap.Uint("user_id", userID))

	// Get user from database
	user, err := u.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("Database error during get user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if user == nil {
		u.logger.Warn("User not found",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user not found")
	}

	u.logger.Info("User retrieved successfully",
		zap.Uint("user_id", userID),
		zap.String("email", user.Email),
	)

	return &dto.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}, nil
}

// DeleteUser menghapus user (soft delete)
func (u *authUsecase) DeleteUser(ctx context.Context, userID uint) (*dto.DeleteUserResponse, error) {
	u.logger.Debug("Delete user", zap.Uint("user_id", userID))

	// Get user first to check if exists
	user, err := u.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("Database error during delete user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if user == nil {
		u.logger.Warn("Delete failed - user not found",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user not found")
	}

	if user.IsDeleted {
		u.logger.Warn("Delete failed - user already deleted",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user already deleted")
	}

	// Mark user as deleted
	if err := u.authRepo.MarkUserAsDeleted(ctx, userID); err != nil {
		u.logger.Error("Failed to delete user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, errors.New("failed to delete user")
	}

	u.logger.Info("User deleted successfully",
		zap.Uint("user_id", userID),
		zap.String("email", user.Email),
	)

	return &dto.DeleteUserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Message: "User deleted successfully",
	}, nil
}

// CheckEmail mengecek apakah email sudah terdaftar
func (u *authUsecase) CheckEmail(ctx context.Context, req dto.CheckEmailRequest) (*dto.CheckEmailResponse, error) {
	u.logger.Debug("Checking email", zap.String("email", req.Email))

	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Database error during email check",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	exists := user != nil && !user.IsDeleted
	message := "email not registered"
	if exists {
		message = "email already registered"
	}

	return &dto.CheckEmailResponse{
		Email:   req.Email,
		Exists:  exists,
		Message: message,
	}, nil
}

// SendOTP mengirim OTP ke email
func (u *authUsecase) SendOTP(ctx context.Context, req dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	u.logger.Debug("Sending OTP",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	// Validasi email terdaftar untuk password reset
	if req.Purpose == "password_reset" {
		user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			u.logger.Error("Database error during email verification",
				zap.String("email", req.Email),
				zap.Error(err),
			)
			return nil, errors.New("database error")
		}

		if user == nil || user.IsDeleted {
			u.logger.Warn("OTP send failed - email not registered",
				zap.String("email", req.Email),
			)
			return nil, errors.New("email not registered")
		}
	}

	// Generate OTP code (6 digit)
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Simpan OTP ke database
	otp := &entity.OTP{
		Email:     req.Email,
		OTPCode:   otpCode,
		Purpose:   req.Purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := u.authRepo.CreateOTP(ctx, otp); err != nil {
		u.logger.Error("Failed to save OTP",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to generate OTP")
	}

	// Kirim OTP via email
	if err := u.emailService.SendOTP(req.Email, otpCode, req.Purpose); err != nil {
		u.logger.Error("Failed to send OTP email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to send OTP email")
	}

	u.logger.Info("OTP sent successfully",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	return &dto.SendOTPResponse{
		Email:   req.Email,
		Message: "OTP has been sent to your email. Valid for 10 minutes.",
	}, nil
}

// ValidateOTP memvalidasi OTP code
func (u *authUsecase) ValidateOTP(ctx context.Context, req dto.ValidateOTPRequest) (*dto.ValidateOTPResponse, error) {
	u.logger.Debug("Validating OTP",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	// Validasi OTP
	isValid, err := u.authRepo.ValidateOTP(ctx, req.Email, req.OTPCode, req.Purpose)
	if err != nil {
		u.logger.Error("Database error during OTP validation",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if !isValid {
		u.logger.Warn("OTP validation failed - invalid or expired OTP",
			zap.String("email", req.Email),
		)
		return nil, errors.New("invalid or expired OTP")
	}

	// Get OTP untuk mendapatkan ID-nya
	otp, err := u.authRepo.GetOTPByEmailAndPurpose(ctx, req.Email, req.Purpose)
	if err != nil {
		u.logger.Error("Database error getting OTP",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	// Mark OTP as used
	if otp != nil {
		if err := u.authRepo.MarkOTPAsUsed(ctx, otp.ID); err != nil {
			u.logger.Error("Failed to mark OTP as used",
				zap.String("email", req.Email),
				zap.Error(err),
			)
		}
	}

	u.logger.Info("OTP validated successfully",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	return &dto.ValidateOTPResponse{
		Valid:   true,
		Message: "OTP is valid",
		Token:   req.OTPCode, // Token untuk reset password (OTP code)
	}, nil
}

// ResetPassword mereset password user
func (u *authUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	u.logger.Debug("Resetting password", zap.String("email", req.Email))

	// Validasi OTP terlebih dahulu
	isValid, err := u.authRepo.ValidateOTP(ctx, req.Email, req.OTPCode, req.Purpose)
	if err != nil {
		u.logger.Error("Database error during OTP validation",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if !isValid {
		u.logger.Warn("Password reset failed - invalid OTP",
			zap.String("email", req.Email),
		)
		return nil, errors.New("invalid or expired OTP")
	}

	// Get user
	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Database error getting user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("database error")
	}

	if user == nil {
		u.logger.Warn("Password reset failed - user not found",
			zap.String("email", req.Email),
		)
		return nil, errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		u.logger.Error("Failed to hash password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to reset password")
	}

	// Update password
	if err := u.authRepo.UpdateUserPassword(ctx, user.ID, hashedPassword); err != nil {
		u.logger.Error("Failed to update user password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errors.New("failed to reset password")
	}

	// Get OTP untuk mark as used
	otp, err := u.authRepo.GetOTPByEmailAndPurpose(ctx, req.Email, req.Purpose)
	if err == nil && otp != nil {
		u.authRepo.MarkOTPAsUsed(ctx, otp.ID)
	}

	u.logger.Info("Password reset successfully",
		zap.String("email", req.Email),
		zap.Uint("user_id", user.ID),
	)

	return &dto.ResetPasswordResponse{
		Email:   req.Email,
		Message: "Password reset successfully. You can now login with your new password.",
	}, nil
}
