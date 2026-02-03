package dto

// LoginRequest merepresentasikan request untuk login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse merepresentasikan response dari login
type LoginResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// CheckEmailRequest merepresentasikan request untuk check email
type CheckEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// CheckEmailResponse merepresentasikan response dari check email
type CheckEmailResponse struct {
	Email   string `json:"email"`
	Exists  bool   `json:"exists"`
	Message string `json:"message"`
}

// SendOTPRequest merepresentasikan request untuk send OTP
type SendOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required,oneof=password_reset email_verification"`
}

// SendOTPResponse merepresentasikan response dari send OTP
type SendOTPResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

// ValidateOTPRequest merepresentasikan request untuk validasi OTP
type ValidateOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required,len=6"`
	Purpose string `json:"purpose" binding:"required,oneof=password_reset email_verification"`
}

// ValidateOTPResponse merepresentasikan response dari validasi OTP
type ValidateOTPResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // Token untuk reset password
}

// ResetPasswordRequest merepresentasikan request untuk reset password
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTPCode     string `json:"otp_code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	Purpose     string `json:"purpose" binding:"required,oneof=password_reset"`
}

// ResetPasswordResponse merepresentasikan response dari reset password
type ResetPasswordResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

// RegisterRequest merepresentasikan request untuk registrasi user baru
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=3"`
	Role     string `json:"role" binding:"omitempty,oneof=customer"` // Optional, auto-set to customer
}

// RegisterResponse merepresentasikan response dari registrasi
type RegisterResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Message string `json:"message"`
}

// GetUserResponse merepresentasikan response dari get user by ID
type GetUserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// DeleteUserResponse merepresentasikan response dari delete user
type DeleteUserResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
