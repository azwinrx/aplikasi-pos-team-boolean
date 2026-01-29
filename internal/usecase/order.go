package usecase

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
)

type OrderUseCase interface {
	GetAllOrders(ctx context.Context) ([]dto.OrderListResponse, error)
	CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (*dto.OrderResponse, error)
	UpdateOrder(ctx context.Context, id uint, req dto.OrderUpdateRequest) error
	DeleteOrder(ctx context.Context, id uint) error
	GetAllTables(ctx context.Context) ([]dto.TableResponse, error)
	GetAllPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error)
	GetAvailableChairs(ctx context.Context) ([]dto.TableResponse, error)
}

type orderUseCase struct {
	orderRepo repository.OrderRepository
	logger    *zap.Logger
}

func NewOrderUseCase(orderRepo repository.OrderRepository, logger *zap.Logger) *orderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (uc *orderUseCase) GetAllOrders(ctx context.Context) ([]dto.OrderListResponse, error) {
	uc.logger.Info("Getting all orders")

	orders, err := uc.orderRepo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all orders", zap.Error(err))
		return nil, err
	}

	var responses []dto.OrderListResponse
	for _, order := range orders {
		responses = append(responses, dto.OrderListResponse{
			ID:           order.ID,
			CustomerName: order.CustomerName,
			TableNumber:  order.Table.Number,
			TotalAmount:  order.TotalAmount,
			Status:       order.Status,
			CreatedAt:    order.CreatedAt,
		})
	}

	uc.logger.Info("Successfully retrieved all orders", zap.Int("count", len(responses)))
	return responses, nil
}

func (uc *orderUseCase) CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (*dto.OrderResponse, error) {
	uc.logger.Info("Creating order",
		zap.String("customer_name", req.CustomerName),
		zap.Uint("table_id", req.TableID))

	order, err := uc.orderRepo.Create(ctx, req)
	if err != nil {
		uc.logger.Error("Failed to create order",
			zap.String("customer_name", req.CustomerName),
			zap.Error(err))
		return nil, err
	}

	// Fetch the full order with preloads
	createdOrder, err := uc.orderRepo.FindByID(ctx, order.ID)
	if err != nil {
		uc.logger.Error("Failed to find created order",
			zap.Uint("id", order.ID),
			zap.Error(err))
		return nil, err
	}

	response := uc.toOrderResponse(*createdOrder)

	uc.logger.Info("Successfully created order",
		zap.Uint("id", order.ID),
		zap.String("customer_name", req.CustomerName))
	return &response, nil
}

func (uc *orderUseCase) UpdateOrder(ctx context.Context, id uint, req dto.OrderUpdateRequest) error {
	uc.logger.Info("Updating order",
		zap.Uint("id", id),
		zap.String("customer_name", req.CustomerName))

	err := uc.orderRepo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("Failed to update order",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	uc.logger.Info("Successfully updated order", zap.Uint("id", id))
	return nil
}

func (uc *orderUseCase) DeleteOrder(ctx context.Context, id uint) error {
	uc.logger.Info("Deleting order", zap.Uint("id", id))

	err := uc.orderRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to delete order",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	uc.logger.Info("Successfully deleted order", zap.Uint("id", id))
	return nil
}

func (uc *orderUseCase) GetAllTables(ctx context.Context) ([]dto.TableResponse, error) {
	uc.logger.Info("Getting all tables")

	tables, err := uc.orderRepo.FindAllTables(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all tables", zap.Error(err))
		return nil, err
	}

	var responses []dto.TableResponse
	for _, table := range tables {
		responses = append(responses, dto.TableResponse{
			ID:       table.ID,
			Number:   table.Number,
			Capacity: table.Capacity,
			Status:   table.Status,
		})
	}

	uc.logger.Info("Successfully retrieved all tables", zap.Int("count", len(responses)))
	return responses, nil
}

func (uc *orderUseCase) GetAllPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error) {
	uc.logger.Info("Getting all payment methods")

	paymentMethods, err := uc.orderRepo.FindAllPaymentMethods(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all payment methods", zap.Error(err))
		return nil, err
	}

	var responses []dto.PaymentMethodResponse
	for _, pm := range paymentMethods {
		responses = append(responses, dto.PaymentMethodResponse{
			ID:   pm.ID,
			Name: pm.Name,
		})
	}

	uc.logger.Info("Successfully retrieved all payment methods", zap.Int("count", len(responses)))
	return responses, nil
}

func (uc *orderUseCase) GetAvailableChairs(ctx context.Context) ([]dto.TableResponse, error) {
	uc.logger.Info("Getting available chairs")

	tables, err := uc.orderRepo.FindAvailableChairs(ctx)
	if err != nil {
		uc.logger.Error("Failed to get available chairs", zap.Error(err))
		return nil, err
	}

	var responses []dto.TableResponse
	for _, table := range tables {
		responses = append(responses, dto.TableResponse{
			ID:       table.ID,
			Number:   table.Number,
			Capacity: table.Capacity,
			Status:   table.Status,
		})
	}

	uc.logger.Info("Successfully retrieved available chairs", zap.Int("count", len(responses)))
	return responses, nil
}

func (uc *orderUseCase) toOrderResponse(order entity.Order) dto.OrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
		})
	}

	return dto.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		TableID:         order.TableID,
		PaymentMethodID: order.PaymentMethodID,
		CustomerName:    order.CustomerName,
		TotalAmount:     order.TotalAmount,
		Tax:             order.Tax,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		Items:           items,
	}
}
