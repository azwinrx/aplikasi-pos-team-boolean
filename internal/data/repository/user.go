package repository

import (
	"context"
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context, f dto.StaffFilterRequest) ([]entity.User, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(ctx context.Context, f dto.StaffFilterRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var totalItems int64

	query := r.db.Model(&entity.User{})

	// Search filter
	if f.Search != "" {
		searchPattern := "%" + strings.ToLower(f.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern)
	}

	// Role filter
	if f.Role != "" {
		query = query.Where("role = ?", f.Role)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortColumn := "name"
	if f.SortBy == "email" {
		sortColumn = "email"
	}

	sortOrder := "asc"
	if f.SortOrder == "desc" {
		sortOrder = "desc"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	offset := (f.Page - 1) * f.Limit
	err := query.Limit(f.Limit).Offset(offset).Find(&users).Error

	return users, totalItems, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
