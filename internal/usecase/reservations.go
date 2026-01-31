package usecase

import (
	"context"
	"errors"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
)

type ReservationsUseCase interface {
	GetAllReservations(ctx context.Context) ([]dto.ReservationListItem, error)
	CreateReservation(ctx context.Context, req dto.ReservationCreateRequest) (*dto.ReservationDetail, error)
	UpdateReservation(ctx context.Context, id uint, req dto.ReservationUpdateRequest) error
	DeleteReservation(ctx context.Context, id uint) error
}

type ReservationUseCase struct {
	repo   repository.ReservationsRepository
	logger *zap.Logger
}

func NewReservationUseCase(repo repository.ReservationsRepository, logger *zap.Logger) *ReservationUseCase {
	return &ReservationUseCase{
		repo:   repo,
		logger: logger,
	}
}

// GetAllReservations returns list of reservations
func (uc *ReservationUseCase) GetAllReservations(ctx context.Context) ([]dto.ReservationListItem, error) {
	results, err := uc.repo.GetAllList(ctx, 0, 1000) // default offset 0, limit 1000
	if err != nil {
		return nil, err
	}
	var list []dto.ReservationListItem
	for _, r := range results {
		list = append(list, dto.ReservationListItem{
			ID:              r.ID,
			CustomerName:    r.CustomerName,
			CustomerPhone:   r.CustomerPhone,
			ReservationDate: r.ReservationDate,
			CheckIn:         r.CheckIn,
			CheckOut:        r.CheckOut,
			TotalPrice:      r.TotalPrice,
		})
	}
	return list, nil
}

// GetReservationByID returns reservation detail by id
func (uc *ReservationUseCase) GetReservationByID(ctx context.Context, id int64) (*dto.ReservationDetail, error) {
	res, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	// TableNumber, PaxNumber, ReservedEnd, DepositFee, TotalPrice perlu join atau query tambahan jika ingin lengkap
	return &dto.ReservationDetail{
		ID:              res.ID,
		CustomerName:    res.CustomerName,
		CustomerPhone:   res.CustomerPhone,
		TableNumber:     "", // butuh join ke table
		PaxNumber:       0,  // butuh join ke table
		ReserveDate:     "",
		ReservationTime: derefTime(res.ReservationTime),
		ReservedEnd:     reservedEnd(res.ReservationTime, 120), // default 2 jam
		DepositFee:      0,                                     // butuh join ke table
		Status:          res.Status,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		TotalPrice:      0, // butuh join ke orders
	}, nil
}

// CreateReservation creates a new reservation from DTO
func (uc *ReservationUseCase) CreateReservation(ctx context.Context, req dto.ReservationCreateRequest) (*dto.ReservationDetail, error) {
	uc.logger.Info("Creating reservation", zap.String("table_number", req.TableNumber), zap.String("customer_name", req.CustomerName))

	// Lookup TableID from TableNumber
	tableID, err := uc.repo.GetTableIDByNumber(ctx, req.TableNumber)
	if err != nil {
		uc.logger.Error("Failed to lookup table ID", zap.String("table_number", req.TableNumber), zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Table ID found", zap.Int64("table_id", tableID))

	// Mapping DTO ke entity
	resTime, err := time.Parse("2006-01-02 15:04", req.ReserveDate+" "+req.ReservationTime)
	if err != nil {
		uc.logger.Error("Failed to parse reservation time", zap.String("reserve_date", req.ReserveDate), zap.String("reservation_time", req.ReservationTime), zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Reservation time parsed", zap.Time("reservation_time", resTime))

	entityReservation := &entity.Reservations{
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		TableID:         tableID,
		ReservationTime: &resTime,
		Status:          req.Status,
	}
	uc.logger.Info("Entity created", zap.Int64("table_id", entityReservation.TableID))
	created, err := uc.repo.Create(ctx, entityReservation)
	if err != nil {
		return nil, err
	}
	return &dto.ReservationDetail{
		ID:              created.ID,
		CustomerName:    created.CustomerName,
		CustomerPhone:   created.CustomerPhone,
		TableNumber:     req.TableNumber,
		PaxNumber:       req.PaxNumber,
		ReserveDate:     req.ReserveDate,
		ReservationTime: *created.ReservationTime,
		ReservedEnd:     created.ReservationTime.Add(time.Duration(req.DurationMinutes) * time.Minute),
		DepositFee:      req.DepositFee,
		Status:          created.Status,
		CreatedAt:       created.CreatedAt,
		UpdatedAt:       created.UpdatedAt,
		TotalPrice:      0, // butuh join ke orders
	}, nil
}

// UpdateReservation updates reservation from DTO
func (uc *ReservationUseCase) UpdateReservation(ctx context.Context, id uint, req dto.ReservationUpdateRequest) error {
	// Mapping TableNumber to TableID - assuming TableNumber is the table number string, need to lookup ID
	// For now, placeholder: assume TableNumber is actually ID as string or something
	// TODO: implement table lookup by number
	var tableID int64 = 0 // placeholder
	if req.TableNumber != "" {
		// lookup tableID by req.TableNumber
		// tableID = getTableIDByNumber(req.TableNumber)
	}
	return uc.repo.Update(ctx, int64(id), tableID, req.Status)
}

// DeleteReservation soft deletes a reservation
func (uc *ReservationUseCase) DeleteReservation(ctx context.Context, id uint) error {
	// Assuming soft delete by setting deleted_at
	return uc.repo.Delete(ctx, int64(id))
}

// Helper untuk dereference *time.Time
func derefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

// Helper untuk reserved end
func reservedEnd(t *time.Time, durationMinutes int) time.Time {
	if t != nil {
		return t.Add(time.Duration(durationMinutes) * time.Minute)
	}
	return time.Time{}
}
