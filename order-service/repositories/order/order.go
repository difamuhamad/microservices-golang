package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	FindAllWithPagination(context.Context, *dto.OrderRequestParam) ([]models.Order, int64, error)
	FindByUUID(context.Context, string) (*models.Order, error)
	FindByUserID(context.Context, string) ([]models.Order, error)
	Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error)
	Update(context.Context, *gorm.DB, *models.Order, uuid.UUID) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

// Find All With Pagination
func (o *OrderRepository) FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error) {
	var (
		orders []models.Order
		sort   string
		total  int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	// Get the data
	err := o.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	// Get the count data
	err = o.db.WithContext(ctx).Model(&models.Order{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, total, nil
}

// Find by UUID
func (o *OrderRepository) FindByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	var order models.Order
	err := o.db.WithContext(ctx).Where("uuid = $1", uuid).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &order, nil
}

// Find by User ID
func (o *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.WithContext(ctx).Where("user_id = $1", userID).Find(&orders).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, nil
}

// Create order code (increment)
func (o *OrderRepository) incrementCode(ctx context.Context) (*string, error) {
	var (
		order  *models.Order
		result string
		today  = time.Now().Format("2006102")
	)

	// check the recent data from db
	err := o.db.WithContext(ctx).Order("id desc").First(&order).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errConstant.ErrSQLError)
		}
	}

	if order.ID != 0 {
		orderCode := order.Code
		splitOrderName, _ := strconv.Atoi(orderCode[4:9])
		code := splitOrderName + 1
		result = fmt.Sprint("ORD-%5d-%s", code, today)
	} else {
		result = fmt.Sprint("ORD-%5d-%s", 1, today)
	}
	// (ORD-00001-2010002) sample output
	return &result, nil
}

// Create
func (o *OrderRepository) Create(ctx context.Context, tx *gorm.DB, param *models.Order) (*models.Order, error) {
	code, err := o.incrementCode(ctx)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UUID:   uuid.New(),
		Code:   *code,
		UserID: param.UserID,
		Amount: param.Amount,
		Date:   param.Date,
		Status: param.Status,
		IsPaid: param.IsPaid,
	}

	err = tx.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return order, nil
}

// Update
func (o *OrderRepository) Update(ctx context.Context, tx *gorm.DB, request *models.Order, uuid uuid.UUID) error {
	err := tx.WithContext(ctx).Model(&models.Order{}).Where("uuid = $1").Updates(request).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
