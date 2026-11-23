package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &OrderRepository{db: db}
}

// CreateWithItems executes an ACID transaction:
// 1. Locks each product row with SELECT ... FOR UPDATE
// 2. Validates product status and available stock
// 3. Decrements stock
// 4. Snapshots product name, price, and subtotal
// 5. Inserts Order, OrderItems, and creates initial Payment record
// 6. Clears the user's cart
func (r *OrderRepository) CreateWithItems(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var calculatedSubtotal int64 = 0
		totalWeightGrams := 0

		for i := range items {
			var product domain.Product
			// Pessimistic lock row for concurrency safety under load
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, items[i].ProductID).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("product ID %d not found: %w", items[i].ProductID, domain.ErrNotFound)
				}
				return err
			}

			if product.Status != "ACTIVE" {
				return fmt.Errorf("product '%s' is not active for sale: %w", product.Name, domain.ErrProductInactive)
			}

			if product.Stock < items[i].Quantity {
				return fmt.Errorf("insufficient stock for '%s' (available: %d, requested: %d): %w",
					product.Name, product.Stock, items[i].Quantity, domain.ErrInsufficientStock)
			}

			// Decrement stock
			newStock := product.Stock - items[i].Quantity
			updateData := map[string]interface{}{"stock": newStock}
			if newStock == 0 {
				updateData["status"] = "OUT_OF_STOCK"
			}
			if err := tx.Model(&product).Updates(updateData).Error; err != nil {
				return fmt.Errorf("failed updating product stock: %w", err)
			}

			// CRITICAL RULE 3: Exact price & weight calculation from locked DB rows
			itemPrice := product.EffectivePrice()
			itemSubtotal := itemPrice * int64(items[i].Quantity)

			items[i].ProductName = product.Name
			items[i].Price = itemPrice
			items[i].Subtotal = itemSubtotal

			calculatedSubtotal += itemSubtotal
			totalWeightGrams += product.Weight * items[i].Quantity
		}

		// Calculate shipping fee & total invoice atomically from locked product weights/prices
		order.ShippingFee = domain.CalculateShippingFee(totalWeightGrams)
		order.TotalAmount = calculatedSubtotal + order.ShippingFee

		// Insert order
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed creating order: %w", err)
		}

		// Insert order items linked to order.ID
		for i := range items {
			items[i].OrderID = order.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return fmt.Errorf("failed creating order items: %w", err)
		}

		// Create initial pending payment record
		initialPayment := domain.Payment{
			OrderID:       order.ID,
			UserID:        order.BuyerID,
			PaymentStatus: "pending",
			GrossAmount:   order.TotalAmount,
		}
		if err := tx.Create(&initialPayment).Error; err != nil {
			return fmt.Errorf("failed creating initial payment record: %w", err)
		}

		// Clear user cart
		var cart domain.Cart
		if err := tx.Where("user_id = ?", order.BuyerID).First(&cart).Error; err == nil {
			tx.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{})
		}

		order.Items = items
		order.Payment = &initialPayment
		return nil
	})
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Buyer").
		Preload("Items.Product").
		Preload("Payment").
		First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByCode(ctx context.Context, code string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Where("order_code = ?", code).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByCodeWithDetails(ctx context.Context, code string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Where("order_code = ?", code).
		Preload("Buyer").
		Preload("Items.Product.Category").
		Preload("Payment").
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetAll(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Order{}).
		Preload("Buyer").
		Preload("Items").
		Preload("Payment")

	if filter.BuyerID != nil && *filter.BuyerID > 0 {
		query = query.Where("buyer_id = ?", *filter.BuyerID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Search != "" {
		term := escapeLikePattern(filter.Search)
		query = query.Where(`order_code ILIKE ? ESCAPE '\' OR recipient_name ILIKE ? ESCAPE '\' OR recipient_phone ILIKE ? ESCAPE '\'`, term, term, term)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status string, trackingNumber string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if trackingNumber != "" {
		updates["tracking_number"] = trackingNumber
	}

	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// RestockCancelledOrder restores product stock if an order expires, cancels, or fails
func (r *OrderRepository) RestockCancelledOrder(ctx context.Context, orderID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, orderID).Error; err != nil {
			return err
		}
		if order.Status == string(domain.OrderStatusCancelled) {
			return nil
		}

		var items []domain.OrderItem
		if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
			return err
		}

		for _, item := range items {
			if err := tx.Model(&domain.Product{}).
				Where("id = ?", item.ProductID).
				Updates(map[string]interface{}{
					"stock":  gorm.Expr("stock + ?", item.Quantity),
					"status": "ACTIVE",
				}).Error; err != nil {
				return err
			}
		}

		return tx.Model(&domain.Order{}).
			Where("id = ?", orderID).
			Update("status", string(domain.OrderStatusCancelled)).Error
	})
}
