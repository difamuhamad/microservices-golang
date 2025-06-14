package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "payment-service/common/error"
	"payment-service/constants"
	errConstant "payment-service/constants/error"
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	FindAllWithPagination(context.Context, *dto.PaymentRequestParam) ([]models.Payment, int64, error)
	FindByUUID(context.Context, string) (*models.Payment, error)
	FindByOrderID(context.Context, string) (*models.Payment, error)
	Create(context.Context, *gorm.DB, *dto.PaymentRequest) (*models.Payment, error)
	Update(context.Context, *gorm.DB, string, *dto.UpdatePaymentRequest) (*models.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

// Get All with Pagination
func (p *PaymentRepository) FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	var (
		payments []models.Payment
		sort     string
		total    int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit

	// Ambil data dengan limit dan offset
	err := p.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&payments).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Hitung total data (tanpa limit & offset)
	err = p.db.
		WithContext(ctx).
		Model(&models.Payment{}).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return payments, total, nil
}

// Find by UUID
func (p *PaymentRepository) FindByUUID(ctx context.Context, uuid string) (*models.Payment, error) {
	var payment models.Payment
	err := p.db.WithContext(ctx).Where("uuid = $1", uuid).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &payment, nil
}

// Find by Order ID
func (p *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := p.db.WithContext(ctx).Where("order_id = $1", orderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &payment, nil
}

// Create Payment
func (p *PaymentRepository) Create(ctx context.Context, tx *gorm.DB, req *dto.PaymentRequest) (*models.Payment, error) {
	status := constants.Initial
	orderID := uuid.MustParse(req.OrderID)
	payment := models.Payment{
		UUID:        uuid.New(),
		OrderID:     orderID,
		Amount:      req.Amount,
		PaymentLink: req.PaymentLink,
		ExpiredAt:   &req.ExpiredAt,
		Description: req.Description,
		Status:      &status,
	}

	// use gorm database transaction (tx)
	err := tx.WithContext(ctx).Create(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &payment, nil
}

// Update
func (p *PaymentRepository) Update(ctx context.Context, tx *gorm.DB, orderID string, req *dto.UpdatePaymentRequest) (*models.Payment, error) {
	payment := models.Payment{
		Status:        req.Status,
		TransactionID: req.TransactionID,
		InvoiceLink:   req.InvoiceLink,
		PaidAt:        req.PaidAt,
		VANumber:      req.VANumber,
		Bank:          req.Bank,
		Acquirer:      req.Acquirer,
	}

	// use gorm database transaction (tx)
	err := tx.WithContext(ctx).Where("order_id = $1", orderID).Updates(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &payment, nil
}
