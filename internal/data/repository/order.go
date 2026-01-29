package repository

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll(ctx context.Context) ([]entity.Order, error)
	FindByID(ctx context.Context, id uint) (*entity.Order, error)
	Create(ctx context.Context, req dto.OrderCreateRequest) (*entity.Order, error)
	Update(ctx context.Context, id uint, req dto.OrderUpdateRequest) error
	Delete(ctx context.Context, id uint) error
	FindAllTables(ctx context.Context) ([]entity.Table, error)
	FindAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error)
	FindAvailableChairs(ctx context.Context) ([]entity.Table, error)
}

type orderRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewOrderRepository(db *gorm.DB, logger *zap.Logger) OrderRepository {
	return &orderRepository{db, logger}
}

func (r *orderRepository) FindAll(ctx context.Context) ([]entity.Order, error) {
	r.logger.Info("Finding all orders")

	var orders []entity.Order
	err := r.db.WithContext(ctx).Preload("Items").Preload("Table").Preload("PaymentMethod").Find(&orders).Error
	if err != nil {
		r.logger.Error("Failed to find all orders", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found all orders", zap.Int("count", len(orders)))
	return orders, nil
}

func (r *orderRepository) FindByID(ctx context.Context, id uint) (*entity.Order, error) {
	r.logger.Info("Finding order by ID", zap.Uint("id", id))

	var order entity.Order
	err := r.db.WithContext(ctx).Preload("Items").Preload("Table").Preload("PaymentMethod").First(&order, id).Error
	if err != nil {
		r.logger.Error("Failed to find order by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found order by ID", zap.Uint("id", order.ID))
	return &order, nil
}

func (r *orderRepository) Create(ctx context.Context, req dto.OrderCreateRequest) (*entity.Order, error) {
	r.logger.Info("Creating new order",
		zap.String("customer_name", req.CustomerName),
		zap.Uint("table_id", req.TableID))

	// Calculate total amount from items
	var totalAmount float64
	var orderItems []entity.OrderItem
	for _, item := range req.Items {
		subtotal := item.Price * float64(item.Quantity)
		totalAmount += subtotal
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  subtotal,
		})
	}
	totalAmount += req.Tax

	order := entity.Order{
		UserID:          req.UserID,
		TableID:         req.TableID,
		PaymentMethodID: req.PaymentMethodID,
		CustomerName:    req.CustomerName,
		TotalAmount:     totalAmount,
		Tax:             req.Tax,
		Status:          "pending",
		Items:           orderItems,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		r.logger.Error("Failed to create order",
			zap.String("customer_name", req.CustomerName),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully created order",
		zap.Uint("id", order.ID),
		zap.String("customer_name", req.CustomerName))
	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, id uint, req dto.OrderUpdateRequest) error {
	r.logger.Info("Updating order",
		zap.Uint("id", id),
		zap.String("customer_name", req.CustomerName))

	// Calculate total amount from items
	var totalAmount float64
	var orderItems []entity.OrderItem
	for _, item := range req.Items {
		subtotal := item.Price * float64(item.Quantity)
		totalAmount += subtotal
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  subtotal,
		})
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order
		if err := tx.Model(&entity.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
			"customer_name":     req.CustomerName,
			"payment_method_id": req.PaymentMethodID,
			"total_amount":      totalAmount,
		}).Error; err != nil {
			return err
		}

		// Delete old items and insert new ones
		if err := tx.Where("order_id = ?", id).Delete(&entity.OrderItem{}).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = id
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.logger.Error("Failed to update order",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully updated order",
		zap.Uint("id", id))
	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting order",
		zap.Uint("id", id))

	// Soft delete
	err := r.db.WithContext(ctx).Delete(&entity.Order{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete order",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted order",
		zap.Uint("id", id))
	return nil
}

func (r *orderRepository) FindAllTables(ctx context.Context) ([]entity.Table, error) {
	r.logger.Info("Finding all tables")

	var tables []entity.Table
	err := r.db.WithContext(ctx).Find(&tables).Error
	if err != nil {
		r.logger.Error("Failed to find all tables", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found all tables", zap.Int("count", len(tables)))
	return tables, nil
}

func (r *orderRepository) FindAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error) {
	r.logger.Info("Finding all payment methods")

	var paymentMethods []entity.PaymentMethod
	err := r.db.WithContext(ctx).Find(&paymentMethods).Error
	if err != nil {
		r.logger.Error("Failed to find all payment methods", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found all payment methods", zap.Int("count", len(paymentMethods)))
	return paymentMethods, nil
}

func (r *orderRepository) FindAvailableChairs(ctx context.Context) ([]entity.Table, error) {
	r.logger.Info("Finding available chairs (tables)")

	var tables []entity.Table
	err := r.db.WithContext(ctx).Where("status = ?", "available").Find(&tables).Error
	if err != nil {
		r.logger.Error("Failed to find available chairs", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found available chairs", zap.Int("count", len(tables)))
	return tables, nil
}
