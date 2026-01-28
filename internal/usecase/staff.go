package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffUseCase interface {
	GetListStaff(ctx context.Context, req dto.StaffFilterRequest) ([]dto.StaffResponse, dto.Pagination, error)
	GetStaffByID(ctx context.Context, id uint) (*dto.StaffResponse, error)
	CreateStaff(ctx context.Context, req dto.StaffCreateRequest) (*dto.StaffResponse, error)
	UpdateStaff(ctx context.Context, id uint, req dto.StaffUpdateRequest) (*dto.StaffResponse, error)
	DeleteStaff(ctx context.Context, id uint) error
}

type staffUseCase struct {
	repo repository.Repository
}

func NewStaffUseCase(repo repository.Repository) StaffUseCase {
	return &staffUseCase{repo}
}

func (s *staffUseCase) GetListStaff(ctx context.Context, req dto.StaffFilterRequest) ([]dto.StaffResponse, dto.Pagination, error) {
	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "name"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Get data from repository
	users, total, err := s.repo.UserRepo.FindAll(ctx, req)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	// Convert to response
	var staffResponses []dto.StaffResponse
	for _, user := range users {
		staffResponses = append(staffResponses, dto.StaffResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
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

func (s *staffUseCase) GetStaffByID(ctx context.Context, id uint) (*dto.StaffResponse, error) {
	user, err := s.repo.UserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, err
	}

	response := &dto.StaffResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *staffUseCase) CreateStaff(ctx context.Context, req dto.StaffCreateRequest) (*dto.StaffResponse, error) {
	// Check if email already exists
	existingUser, err := s.repo.UserRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Generate random password
	password := generateRandomPassword(12)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	// Save to database
	if err := s.repo.UserRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// TODO: Send email with password
	// Implementasi pengiriman email dengan password yang digenerate
	fmt.Printf("Generated password for %s: %s\n", user.Email, password)

	response := &dto.StaffResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *staffUseCase) UpdateStaff(ctx context.Context, id uint, req dto.StaffUpdateRequest) (*dto.StaffResponse, error) {
	// Get existing user
	user, err := s.repo.UserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, err
	}

	// Check if email already exists (exclude current user)
	if req.Email != user.Email {
		existingUser, err := s.repo.UserRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
	}

	// Update user data
	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	// Save to database
	if err := s.repo.UserRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := &dto.StaffResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *staffUseCase) DeleteStaff(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := s.repo.UserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("staff not found")
		}
		return err
	}

	// Soft delete
	if err := s.repo.UserRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// generateRandomPassword generates a random password with specified length
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}
	return string(password)
}
