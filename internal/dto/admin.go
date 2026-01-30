package dto

import "time"

// AdminResponse adalah response untuk data admin
type AdminResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListAdminResponse adalah response untuk list admin
type ListAdminResponse struct {
	Data       []AdminResponse `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// EditAdminAccessRequest adalah request untuk edit akses admin
type EditAdminAccessRequest struct {
	Role   string `json:"role" binding:"required,oneof=admin superadmin user"`
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// EditAdminAccessResponse adalah response untuk edit akses admin
type EditAdminAccessResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreateAdminRequest adalah request untuk membuat admin baru
type CreateAdminRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=3,max=100"`
	Role  string `json:"role" binding:"required,oneof=admin superadmin"`
}

// UpdateUserProfileRequest adalah request untuk update profil user
type UpdateUserProfileRequest struct {
	Name     string `json:"name" binding:"max=100"`
	Password string `json:"password" binding:"min=6"`
}

// UserProfileResponse adalah response untuk profil user
type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Message   string    `json:"message,omitempty"`
}

// LogoutRequest adalah request untuk logout
type LogoutRequest struct {
	UserID uint `json:"user_id"`
}

// LogoutResponse adalah response untuk logout
type LogoutResponse struct {
	Message string `json:"message"`
}
