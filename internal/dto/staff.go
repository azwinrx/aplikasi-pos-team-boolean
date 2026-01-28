package dto

// StaffFilterRequest untuk filter staff list
type StaffFilterRequest struct {
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
	Search    string `json:"search" form:"search"`         // Cari berdasarkan full_name, email, phone_number
	SortBy    string `json:"sort_by" form:"sort_by"`       // full_name, email, salary, date_of_birth, created_at
	SortOrder string `json:"sort_order" form:"sort_order"` // asc, desc
	Role      string `json:"role" form:"role"`             // Filter by role (admin, manager, cashier, staff, supervisor)
}

// StaffCreateRequest untuk create staff baru - Sesuai UI Form
type StaffCreateRequest struct {
	FullName          string  `json:"full_name" binding:"required,min=3,max=100"`
	Email             string  `json:"email" binding:"required,email"`
	Role              string  `json:"role" binding:"required,oneof=admin manager cashier staff supervisor"`
	PhoneNumber       string  `json:"phone_number"`
	Salary            float64 `json:"salary" binding:"min=0"`
	DateOfBirth       string  `json:"date_of_birth"`      // Format: YYYY-MM-DD
	ShiftStartTiming  string  `json:"shift_start_timing"` // Format: HH:MM:SS atau HH:MM
	ShiftEndTiming    string  `json:"shift_end_timing"`   // Format: HH:MM:SS atau HH:MM
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details"`
}

// StaffUpdateRequest untuk update staff - Sesuai UI Form
type StaffUpdateRequest struct {
	FullName          string  `json:"full_name" binding:"required,min=3,max=100"`
	Email             string  `json:"email" binding:"required,email"`
	Role              string  `json:"role" binding:"required,oneof=admin manager cashier staff supervisor"`
	PhoneNumber       string  `json:"phone_number"`
	Salary            float64 `json:"salary" binding:"min=0"`
	DateOfBirth       string  `json:"date_of_birth"`      // Format: YYYY-MM-DD
	ShiftStartTiming  string  `json:"shift_start_timing"` // Format: HH:MM:SS atau HH:MM
	ShiftEndTiming    string  `json:"shift_end_timing"`   // Format: HH:MM:SS atau HH:MM
	Address           string  `json:"address"`
	AdditionalDetails string  `json:"additional_details"`
}

// StaffResponse untuk response staff - Sesuai UI Display
type StaffResponse struct {
	ID                uint    `json:"id"`
	FullName          string  `json:"full_name"`
	Email             string  `json:"email"`
	Role              string  `json:"role"`
	PhoneNumber       string  `json:"phone_number,omitempty"`
	Salary            float64 `json:"salary"`
	DateOfBirth       string  `json:"date_of_birth,omitempty"`
	ShiftStartTiming  string  `json:"shift_start_timing,omitempty"`
	ShiftEndTiming    string  `json:"shift_end_timing,omitempty"`
	Address           string  `json:"address,omitempty"`
	AdditionalDetails string  `json:"additional_details,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// StaffListResponse untuk response list staff yang sederhana
type StaffListResponse struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   string  `json:"phone"`
	Age     int     `json:"age"`
	Salary  float64 `json:"salary"`
	Timings string  `json:"timings"` // Format: "9am to 6pm"
}

// StaffDetailResponse untuk response detail staff yang sederhana
type StaffDetailResponse struct {
	FullName         string  `json:"full_name"`
	Email            string  `json:"email"`
	PhoneNumber      string  `json:"phone_number"`
	DateOfBirth      string  `json:"date_of_birth"`
	Address          string  `json:"address"`
	Role             string  `json:"role"`
	Salary           float64 `json:"salary"`
	ShiftStartTiming string  `json:"shift_start_timing"`
	ShiftEndTiming   string  `json:"shift_end_timing"`
}
