package entity

import (
	"time"

	"gorm.io/gorm"
)

// Staff merepresentasikan tabel staff di database (data pegawai lengkap)
type Staff struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName          string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email             string         `gorm:"type:varchar(100);unique;not null;index" json:"email"`
	Role              string         `gorm:"type:varchar(20);not null;default:'staff';index" json:"role"`
	PhoneNumber       string         `gorm:"type:varchar(20)" json:"phone_number"`
	Salary            float64        `gorm:"type:decimal(15,2);default:0" json:"salary"`
	DateOfBirth       *time.Time     `gorm:"type:date" json:"date_of_birth,omitempty"`
	ShiftStartTiming  string         `gorm:"type:varchar(10)" json:"shift_start_timing,omitempty"`
	ShiftEndTiming    string         `gorm:"type:varchar(10)" json:"shift_end_timing,omitempty"`
	Address           string         `gorm:"type:text" json:"address"`
	AdditionalDetails string         `gorm:"type:text" json:"additional_details"`
	CreatedAt         time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (Staff) TableName() string {
	return "staff"
}

// BeforeCreate hook untuk validasi sebelum create
func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.Role == "" {
		s.Role = "staff"
	}
	return nil
}

// BeforeUpdate hook untuk update timestamp
func (s *Staff) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
