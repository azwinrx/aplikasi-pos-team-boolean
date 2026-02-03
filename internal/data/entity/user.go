package entity

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan tabel users di database untuk authentication
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(100);unique;not null;index" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Role      string         `gorm:"type:varchar(20);not null;default:'customer';index" json:"role"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, inactive
	IsDeleted bool           `gorm:"default:false;index" json:"is_deleted"`                    // Untuk melacak user yang dihapus
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (User) TableName() string {
	return "users"
}

// OTP merepresentasikan tabel otps di database untuk menyimpan OTP
type OTP struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(100);not null;index" json:"email"`
	OTPCode   string    `gorm:"type:varchar(10);not null" json:"otp_code"`
	Purpose   string    `gorm:"type:varchar(50);not null;default:'password_reset'"` // password_reset, email_verification
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName override nama tabel
func (OTP) TableName() string {
	return "otps"
}
