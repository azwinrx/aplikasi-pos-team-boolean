package repository

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuthRepository mendefinisikan interface untuk auth operations
type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, id uint, hashedPassword string) error
	MarkUserAsDeleted(ctx context.Context, id uint) error

	// OTP operations
	CreateOTP(ctx context.Context, otp *entity.OTP) error
	GetOTPByEmailAndPurpose(ctx context.Context, email, purpose string) (*entity.OTP, error)
	ValidateOTP(ctx context.Context, email, otpCode, purpose string) (bool, error)
	MarkOTPAsUsed(ctx context.Context, otpID uint) error
	DeleteExpiredOTPs(ctx context.Context) error
}

// authRepository implementasi dari AuthRepository interface
type authRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewAuthRepository membuat instance baru dari authRepository
func NewAuthRepository(db *gorm.DB, logger *zap.Logger) AuthRepository {
	return &authRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser membuat user baru di database
func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("Failed to create user",
			zap.String("email", user.Email),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("User created successfully",
		zap.String("email", user.Email),
		zap.Uint("id", user.ID),
	)

	return nil
}

// GetUserByEmail mengambil user berdasarkan email
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	if err := r.db.WithContext(ctx).Where("email = ? AND is_deleted = ?", email, false).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("User not found",
				zap.String("email", email),
			)
			return nil, nil // Return nil untuk user tidak ditemukan
		}

		r.logger.Error("Failed to get user by email",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, err
	}

	return &user, nil
}

// GetUserByID mengambil user berdasarkan ID
func (r *authRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User

	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("User not found",
				zap.Uint("id", id),
			)
			return nil, nil
		}

		r.logger.Error("Failed to get user by ID",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return &user, nil
}

// UpdateUserPassword mengupdate password user
func (r *authRepository) UpdateUserPassword(ctx context.Context, id uint, hashedPassword string) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		r.logger.Error("Failed to update user password",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("User password updated successfully",
		zap.Uint("id", id),
	)

	return nil
}

// MarkUserAsDeleted menandai user sebagai deleted
func (r *authRepository) MarkUserAsDeleted(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_deleted", true).Error; err != nil {
		r.logger.Error("Failed to mark user as deleted",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("User marked as deleted",
		zap.Uint("id", id),
	)

	return nil
}

// CreateOTP membuat OTP baru di database
func (r *authRepository) CreateOTP(ctx context.Context, otp *entity.OTP) error {
	// Set expiry time ke 10 menit dari sekarang jika belum di-set
	if otp.ExpiresAt.IsZero() {
		otp.ExpiresAt = time.Now().Add(10 * time.Minute)
	}

	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		r.logger.Error("Failed to create OTP",
			zap.String("email", otp.Email),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("OTP created successfully",
		zap.String("email", otp.Email),
		zap.String("purpose", otp.Purpose),
	)

	return nil
}

// GetOTPByEmailAndPurpose mengambil OTP berdasarkan email dan purpose yang tidak expired dan belum digunakan
func (r *authRepository) GetOTPByEmailAndPurpose(ctx context.Context, email, purpose string) (*entity.OTP, error) {
	var otp entity.OTP

	if err := r.db.WithContext(ctx).
		Where("email = ? AND purpose = ? AND is_used = ? AND expires_at > ?", email, purpose, false, time.Now()).
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("OTP not found",
				zap.String("email", email),
				zap.String("purpose", purpose),
			)
			return nil, nil
		}

		r.logger.Error("Failed to get OTP",
			zap.String("email", email),
			zap.String("purpose", purpose),
			zap.Error(err),
		)
		return nil, err
	}

	return &otp, nil
}

// ValidateOTP memvalidasi OTP code
func (r *authRepository) ValidateOTP(ctx context.Context, email, otpCode, purpose string) (bool, error) {
	otp, err := r.GetOTPByEmailAndPurpose(ctx, email, purpose)
	if err != nil {
		return false, err
	}

	if otp == nil {
		r.logger.Warn("OTP not found or expired",
			zap.String("email", email),
			zap.String("purpose", purpose),
		)
		return false, nil
	}

	// Validasi OTP code
	if otp.OTPCode != otpCode {
		r.logger.Warn("Invalid OTP code",
			zap.String("email", email),
			zap.String("provided_code", otpCode),
		)
		return false, nil
	}

	return true, nil
}

// MarkOTPAsUsed menandai OTP sebagai sudah digunakan
func (r *authRepository) MarkOTPAsUsed(ctx context.Context, otpID uint) error {
	if err := r.db.WithContext(ctx).Model(&entity.OTP{}).Where("id = ?", otpID).Update("is_used", true).Error; err != nil {
		r.logger.Error("Failed to mark OTP as used",
			zap.Uint("otp_id", otpID),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("OTP marked as used",
		zap.Uint("otp_id", otpID),
	)

	return nil
}

// DeleteExpiredOTPs menghapus OTP yang sudah expired
func (r *authRepository) DeleteExpiredOTPs(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.OTP{}).Error; err != nil {
		r.logger.Error("Failed to delete expired OTPs",
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Expired OTPs deleted successfully")

	return nil
}
