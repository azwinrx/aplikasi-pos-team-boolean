package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StaffUseCase interface {
	GetListStaff(ctx context.Context, req dto.StaffFilterRequest) ([]dto.StaffListResponse, dto.Pagination, error)
	GetStaffByID(ctx context.Context, id uint) (*dto.StaffDetailResponse, error)
	GetStaffByEmail(ctx context.Context, email string) (*dto.StaffResponse, error)
	CreateStaff(ctx context.Context, req dto.StaffCreateRequest) (*dto.StaffResponse, error)
	UpdateStaff(ctx context.Context, id uint, req dto.StaffUpdateRequest) (*dto.StaffResponse, error)
	DeleteStaff(ctx context.Context, id uint) error
}

type staffUseCase struct {
	staffRepo repository.StaffRepository
	logger    *zap.Logger
}

func NewStaffUseCase(staffRepo repository.StaffRepository, logger *zap.Logger) *staffUseCase {
	return &staffUseCase{
		staffRepo: staffRepo,
		logger:    logger,
	}
}

func (s *staffUseCase) GetListStaff(ctx context.Context, req dto.StaffFilterRequest) ([]dto.StaffListResponse, dto.Pagination, error) {
	s.logger.Debug("Fetching staff list",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit),
	)

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "full_name"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Get data from repository
	staffList, total, err := s.staffRepo.FindAll(ctx, req)
	if err != nil {
		s.logger.Error("Failed to fetch staff list from database", zap.Error(err))
		return nil, dto.Pagination{}, err
	}

	s.logger.Info("Staff list fetched successfully",
		zap.Int64("total_items", total),
		zap.Int("returned_items", len(staffList)),
		zap.Int("page", req.Page),
	)

	// Convert to list response
	var staffResponses []dto.StaffListResponse
	for _, staff := range staffList {
		staffResponses = append(staffResponses, s.toStaffListResponse(&staff))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := dto.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		TotalItems: int(total),
	}

	return staffResponses, pagination, nil
}

func (s *staffUseCase) GetStaffByID(ctx context.Context, id uint) (*dto.StaffDetailResponse, error) {
	s.logger.Debug("Fetching staff by ID", zap.Uint("id", id))

	staff, err := s.staffRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Staff not found", zap.Uint("id", id))
			return nil, errors.New("staff not found")
		}
		s.logger.Error("Failed to fetch staff by ID", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	s.logger.Info("Staff fetched successfully", zap.Uint("id", id), zap.String("full_name", staff.FullName))

	response := s.toStaffDetailResponse(staff)
	return &response, nil
}

func (s *staffUseCase) GetStaffByEmail(ctx context.Context, email string) (*dto.StaffResponse, error) {
	s.logger.Debug("Fetching staff by email", zap.String("email", email))

	staff, err := s.staffRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Staff not found", zap.String("email", email))
			return nil, errors.New("staff not found")
		}
		s.logger.Error("Failed to fetch staff by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	s.logger.Info("Staff fetched successfully by email", zap.Uint("id", staff.ID), zap.String("email", staff.Email))

	response := s.toStaffResponse(staff)
	return &response, nil
}

func (s *staffUseCase) CreateStaff(ctx context.Context, req dto.StaffCreateRequest) (*dto.StaffResponse, error) {
	s.logger.Info("Creating new staff",
		zap.String("full_name", req.FullName),
		zap.String("email", req.Email),
		zap.String("role", req.Role),
	)

	// Validate required fields
	if req.FullName == "" {
		s.logger.Warn("Validation failed: full name is required")
		return nil, errors.New("full name is required")
	}
	if req.Email == "" {
		s.logger.Warn("Validation failed: email is required")
		return nil, errors.New("email is required")
	}
	if req.Role == "" {
		s.logger.Warn("Validation failed: role is required")
		return nil, errors.New("role is required")
	}

	// Check if email already exists
	existingStaff, err := s.staffRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingStaff != nil {
		s.logger.Warn("Email already exists", zap.String("email", req.Email))
		return nil, errors.New("email already exists")
	}

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			s.logger.Warn("Invalid date of birth format", zap.Error(err), zap.String("date", req.DateOfBirth))
			return nil, fmt.Errorf("invalid date of birth format. Use YYYY-MM-DD: %w", err)
		}
		dateOfBirth = &parsedDate
	}

	// Shift timing is stored as string, no parsing needed

	// Create staff entity
	staff := &entity.Staff{
		FullName:          req.FullName,
		Email:             req.Email,
		Role:              req.Role,
		PhoneNumber:       req.PhoneNumber,
		Salary:            req.Salary,
		DateOfBirth:       dateOfBirth,
		ShiftStartTiming:  req.ShiftStartTiming,
		ShiftEndTiming:    req.ShiftEndTiming,
		Address:           req.Address,
		AdditionalDetails: req.AdditionalDetails,
	}

	// Save to database
	if err := s.staffRepo.Create(ctx, staff); err != nil {
		s.logger.Error("Failed to create staff in database",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		return nil, err
	}

	s.logger.Info("Staff created successfully",
		zap.Uint("id", staff.ID),
		zap.String("full_name", staff.FullName),
		zap.String("email", staff.Email),
	)

	response := s.toStaffResponse(staff)
	return &response, nil
}

func (s *staffUseCase) UpdateStaff(ctx context.Context, id uint, req dto.StaffUpdateRequest) (*dto.StaffResponse, error) {
	s.logger.Info("Updating staff",
		zap.Uint("id", id),
		zap.String("full_name", req.FullName),
		zap.String("email", req.Email),
		zap.String("role", req.Role),
	)

	// Get existing staff
	staff, err := s.staffRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Staff not found for update", zap.Uint("id", id))
			return nil, errors.New("staff not found")
		}
		s.logger.Error("Failed to fetch staff for update", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	// Check if email already exists (exclude current staff)
	if req.Email != staff.Email {
		existingStaff, err := s.staffRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingStaff != nil && existingStaff.ID != id {
			s.logger.Warn("Email already exists for another staff",
				zap.String("email", req.Email),
				zap.Uint("existing_staff_id", existingStaff.ID),
			)
			return nil, errors.New("email already exists")
		}
	}

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			s.logger.Warn("Invalid date of birth format", zap.Error(err), zap.String("date", req.DateOfBirth))
			return nil, fmt.Errorf("invalid date of birth format. Use YYYY-MM-DD: %w", err)
		}
		dateOfBirth = &parsedDate
	}

	// Shift timing is stored as string, no parsing needed

	// Update staff data
	staff.FullName = req.FullName
	staff.Email = req.Email
	staff.Role = req.Role
	staff.PhoneNumber = req.PhoneNumber
	staff.Salary = req.Salary
	staff.DateOfBirth = dateOfBirth
	staff.ShiftStartTiming = req.ShiftStartTiming
	staff.ShiftEndTiming = req.ShiftEndTiming
	staff.Address = req.Address
	staff.AdditionalDetails = req.AdditionalDetails

	// Save to database
	if err := s.staffRepo.Update(ctx, staff); err != nil {
		s.logger.Error("Failed to update staff in database",
			zap.Error(err),
			zap.Uint("id", id),
		)
		return nil, err
	}

	s.logger.Info("Staff updated successfully",
		zap.Uint("id", staff.ID),
		zap.String("full_name", staff.FullName),
	)

	response := s.toStaffResponse(staff)
	return &response, nil
}

func (s *staffUseCase) DeleteStaff(ctx context.Context, id uint) error {
	s.logger.Info("Deleting staff", zap.Uint("id", id))

	// Check if staff exists
	staff, err := s.staffRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Staff not found for deletion", zap.Uint("id", id))
			return errors.New("staff not found")
		}
		s.logger.Error("Failed to fetch staff for deletion", zap.Error(err), zap.Uint("id", id))
		return err
	}

	// Soft delete
	if err := s.staffRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete staff from database",
			zap.Error(err),
			zap.Uint("id", id),
		)
		return err
	}

	s.logger.Info("Staff deleted successfully",
		zap.Uint("id", id),
		zap.String("full_name", staff.FullName),
	)

	return nil
}

// toStaffListResponse converts entity to list response DTO (simplified)
func (s *staffUseCase) toStaffListResponse(staff *entity.Staff) dto.StaffListResponse {
	response := dto.StaffListResponse{
		ID:      staff.ID,
		Name:    staff.FullName,
		Email:   staff.Email,
		Phone:   staff.PhoneNumber,
		Salary:  staff.Salary,
		Timings: "N/A", // Default value
	}

	// Calculate age from date of birth
	if staff.DateOfBirth != nil {
		now := time.Now()
		age := now.Year() - staff.DateOfBirth.Year()
		if now.YearDay() < staff.DateOfBirth.YearDay() {
			age--
		}
		response.Age = age
	}

	// Format timings as "9am to 6pm" from string format "HH:MM:SS" or "HH:MM"
	if staff.ShiftStartTiming != "" && staff.ShiftEndTiming != "" {
		// Parse start time
		startTime, err := time.Parse("15:04:05", staff.ShiftStartTiming)
		if err != nil {
			// Try without seconds
			startTime, err = time.Parse("15:04", staff.ShiftStartTiming)
		}

		// Parse end time
		endTime, err2 := time.Parse("15:04:05", staff.ShiftEndTiming)
		if err2 != nil {
			// Try without seconds
			endTime, err2 = time.Parse("15:04", staff.ShiftEndTiming)
		}

		// If both parsed successfully, format as "9am to 6pm"
		if err == nil && err2 == nil {
			startHour := startTime.Hour()
			endHour := endTime.Hour()

			startAmPm := "am"
			endAmPm := "am"

			if startHour >= 12 {
				startAmPm = "pm"
				if startHour > 12 {
					startHour -= 12
				}
			}
			if startHour == 0 {
				startHour = 12
			}

			if endHour >= 12 {
				endAmPm = "pm"
				if endHour > 12 {
					endHour -= 12
				}
			}
			if endHour == 0 {
				endHour = 12
			}

			response.Timings = fmt.Sprintf("%d%s to %d%s", startHour, startAmPm, endHour, endAmPm)
		}
	}

	return response
}

// toStaffResponse converts entity to response DTO
func (s *staffUseCase) toStaffResponse(staff *entity.Staff) dto.StaffResponse {
	response := dto.StaffResponse{
		ID:                staff.ID,
		FullName:          staff.FullName,
		Email:             staff.Email,
		Role:              staff.Role,
		PhoneNumber:       staff.PhoneNumber,
		Salary:            staff.Salary,
		Address:           staff.Address,
		AdditionalDetails: staff.AdditionalDetails,
		CreatedAt:         staff.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         staff.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if staff.DateOfBirth != nil {
		response.DateOfBirth = staff.DateOfBirth.Format("2006-01-02")
	}

	response.ShiftStartTiming = staff.ShiftStartTiming
	response.ShiftEndTiming = staff.ShiftEndTiming

	return response
}

// toStaffDetailResponse converts entity to detail response DTO (simplified for GetByID)
func (s *staffUseCase) toStaffDetailResponse(staff *entity.Staff) dto.StaffDetailResponse {
	response := dto.StaffDetailResponse{
		FullName:         staff.FullName,
		Email:            staff.Email,
		PhoneNumber:      staff.PhoneNumber,
		Address:          staff.Address,
		Role:             staff.Role,
		Salary:           staff.Salary,
		ShiftStartTiming: staff.ShiftStartTiming,
		ShiftEndTiming:   staff.ShiftEndTiming,
	}

	// Format date of birth
	if staff.DateOfBirth != nil {
		response.DateOfBirth = staff.DateOfBirth.Format("2006-01-02")
	}

	return response
}
