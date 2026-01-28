package repository

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"

	"gorm.io/gorm"
)

type StaffRepository interface {
	FindAll(ctx context.Context, f dto.StaffFilterRequest) ([]entity.Staff, int64, error)
	Create(ctx context.Context, staff *entity.Staff) error
	Update(ctx context.Context, staff *entity.Staff) error
	Detail(ctx context.Context, id uint) (*entity.Staff, error)
	Delete(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*entity.Staff, error)
}

type staffRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewStaffRepository(db *gorm.DB, logger *zap.Logger) StaffRepository {
	return &staffRepository{db, logger}
}

func (r *staffRepository) FindAll(ctx context.Context, f dto.StaffFilterRequest) ([]entity.Staff, int64, error) {
	r.logger.Info("Finding all staff",
		zap.Int("page", f.Page),
		zap.Int("limit", f.Limit))

	var staffList []entity.Staff
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Staff{})

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		r.logger.Error("Failed to count staff items",
			zap.Error(err))
		return nil, 0, err
	}

	// Sorting - default by full_name
	sortColumn := "full_name"
	if f.SortBy != "" {
		switch f.SortBy {
		case "email":
			sortColumn = "email"
		case "salary":
			sortColumn = "salary"
		case "date_of_birth":
			sortColumn = "date_of_birth"
		case "created_at":
			sortColumn = "created_at"
		default:
			sortColumn = "full_name"
		}
	}

	sortOrder := "asc"
	if f.SortOrder == "desc" {
		sortOrder = "desc"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}

	offset := (f.Page - 1) * f.Limit

	// Select only needed fields: id, full_name, email, phone_number, date_of_birth, salary, shift_start_timing, shift_end_timing
	err := query.Select("id", "full_name", "email", "phone_number", "date_of_birth", "salary", "shift_start_timing", "shift_end_timing").
		Limit(f.Limit).
		Offset(offset).
		Find(&staffList).Error

	if err != nil {
		r.logger.Error("Failed to find all staff",
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Successfully found all staff",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(staffList)))

	return staffList, totalItems, err
}

func (r *staffRepository) Create(ctx context.Context, staff *entity.Staff) error {
	r.logger.Info("Creating new staff",
		zap.String("full_name", staff.FullName),
		zap.String("email", staff.Email),
		zap.String("role", staff.Role))

	err := r.db.WithContext(ctx).Create(staff).Error
	if err != nil {
		r.logger.Error("Failed to create staff",
			zap.String("full_name", staff.FullName),
			zap.String("email", staff.Email),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully created staff",
		zap.Uint("id", staff.ID),
		zap.String("full_name", staff.FullName))
	return err
}

func (r *staffRepository) Update(ctx context.Context, staff *entity.Staff) error {
	r.logger.Info("Updating staff",
		zap.Uint("id", staff.ID),
		zap.String("full_name", staff.FullName))

	err := r.db.WithContext(ctx).Save(staff).Error
	if err != nil {
		r.logger.Error("Failed to update staff",
			zap.Uint("id", staff.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully updated staff",
		zap.Uint("id", staff.ID))
	return err
}

func (r *staffRepository) Detail(ctx context.Context, id uint) (*entity.Staff, error) {
	r.logger.Info("Getting staff detail",
		zap.Uint("id", id))

	var staff entity.Staff
	err := r.db.WithContext(ctx).First(&staff, id).Error
	if err != nil {
		r.logger.Error("Failed to get staff detail",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully retrieved staff detail",
		zap.Uint("id", id),
		zap.String("full_name", staff.FullName))
	return &staff, err
}

func (r *staffRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting staff",
		zap.Uint("id", id))

	// Soft delete
	err := r.db.WithContext(ctx).Delete(&entity.Staff{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete staff",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted staff",
		zap.Uint("id", id))
	return err
}

func (r *staffRepository) FindByEmail(ctx context.Context, email string) (*entity.Staff, error) {
	r.logger.Info("Finding staff by email",
		zap.String("email", email))

	var staff entity.Staff
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&staff).Error
	if err != nil {
		r.logger.Error("Failed to find staff by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found staff by email",
		zap.Uint("id", staff.ID),
		zap.String("email", staff.Email))
	return &staff, err
}
