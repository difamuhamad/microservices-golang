package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrderFieldRepository struct {
	db *gorm.DB
}

type IOrderFieldRepository interface {
	FindByOrderID(context.Context, uint) ([]models.OrderField, error)
	Create(context.Context, *gorm.DB, []models.OrderField) error
}

func NewOrderFieldRepository(db *gorm.DB) IOrderFieldRepository {
	return &OrderFieldRepository{db: db}
}

func (o *OrderFieldRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderField, error) {
	var orderFields []models.OrderField
	err := o.db.WithContext(ctx).Where("order_id = $1", orderID).Find(&orderFields).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return orderFields, nil
}

func (o *OrderFieldRepository) Create(ctx context.Context, tx *gorm.DB, request []models.OrderField) error {
	err := tx.WithContext(ctx).Create(&request).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
