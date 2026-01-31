package repository

import (
	"context"
	"errors"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationsRepository interface {
	GetAllList(ctx context.Context, offset, limit int) ([]ReservationListResult, error)
	GetById(ctx context.Context, id int64) (*entity.Reservations, error)
	Create(ctx context.Context, reservation *entity.Reservations) (*entity.Reservations, error)
	Update(ctx context.Context, id int64, newTableID int64, cancelStatus string) error
	Delete(ctx context.Context, id int64) error
	GetTableIDByNumber(ctx context.Context, tableNumber string) (int64, error)
}

type reservationsRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewReservationsRepository(db *gorm.DB, logger *zap.Logger) ReservationsRepository {
	return &reservationsRepository{
		db:     db,
		logger: logger,
	}
}

// GetAllList: List all reservations with pagination (raw entity, mapping ke DTO di usecase)
type ReservationListResult struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	ReservationDate string    `json:"reservation_date"`
	CheckIn         time.Time `json:"check_in"`
	CheckOut        time.Time `json:"check_out"`
	TotalPrice      float64   `json:"total_price"`
}

func (r *reservationsRepository) GetAllList(ctx context.Context, offset, limit int) ([]ReservationListResult, error) {
	var results []ReservationListResult
	// Join ke orders untuk ambil total price, join ke reservations untuk data reservasi
	// Diasumsikan orders.reservation_id = reservations.id (jika tidak, sesuaikan relasi)
	query := `
		SELECT
			reservations.id,
			reservations.customer_name,
			reservations.customer_phone,
			COALESCE(TO_CHAR(reservations.reservation_time, 'YYYY-MM-DD'), '') as reservation_date,
			reservations.reservation_time as check_in,
			(reservations.reservation_time + INTERVAL '2 hour') as check_out,
			COALESCE(SUM(orders.total_amount), 0) as total_price
		FROM reservations
		LEFT JOIN orders ON orders.table_id = reservations.table_id
			AND DATE(orders.created_at) = DATE(reservations.reservation_time)
			AND orders.deleted_at IS NULL
		WHERE reservations.deleted_at IS NULL
		GROUP BY reservations.id, reservations.customer_name, reservations.customer_phone, reservations.reservation_time
		ORDER BY reservations.reservation_time DESC
		OFFSET ? LIMIT ?
	`
	if err := r.db.WithContext(ctx).Raw(query, offset, limit).Scan(&results).Error; err != nil {
		r.logger.Error("failed to get reservation list", zap.Error(err))
		return nil, err
	}
	return results, nil
}

// GetById: Get reservation detail by id (raw entity)
func (r *reservationsRepository) GetById(ctx context.Context, id int64) (*entity.Reservations, error) {
	var res entity.Reservations
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("failed to get reservation by id", zap.Error(err))
		return nil, err
	}
	return &res, nil
}

// Create: Add new reservation (raw entity, mapping dari DTO ke entity di usecase)
func (r *reservationsRepository) Create(ctx context.Context, reservation *entity.Reservations) (*entity.Reservations, error) {
	if reservation == nil {
		return nil, errors.New("reservation is nil")
	}
	if reservation.TableID == 0 {
		return nil, errors.New("table_id is required")
	}
	if reservation.ReservationTime == nil {
		return nil, errors.New("reservation_time is required")
	}
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		r.logger.Error("failed to create reservation", zap.Error(err))
		return nil, err
	}
	return reservation, nil
}

// Update: hanya update table_id dan status pembatalan
func (r *reservationsRepository) Update(ctx context.Context, id int64, newTableID int64, cancelStatus string) error {
	updateData := map[string]interface{}{}
	if newTableID != 0 {
		updateData["table_id"] = newTableID
	}
	if cancelStatus != "" {
		updateData["status"] = cancelStatus
	}
	if len(updateData) == 0 {
		return errors.New("no update data provided")
	}
	if err := r.db.WithContext(ctx).
		Model(&entity.Reservations{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updateData).Error; err != nil {
		r.logger.Error("failed to update reservation", zap.Error(err))
		return err
	}
	return nil
}

// Delete: soft delete reservation by setting deleted_at
func (r *reservationsRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Reservations{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error; err != nil {
		r.logger.Error("failed to delete reservation", zap.Error(err))
		return err
	}
	return nil
}

// GetTableIDByNumber: Get table ID by table number
func (r *reservationsRepository) GetTableIDByNumber(ctx context.Context, tableNumber string) (int64, error) {
	var table entity.Table
	err := r.db.WithContext(ctx).
		Where("number = ? AND deleted_at IS NULL", tableNumber).
		First(&table).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("table not found")
		}
		r.logger.Error("failed to get table by number", zap.Error(err))
		return 0, err
	}
	return int64(table.ID), nil
}
