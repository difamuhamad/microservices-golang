package repositories

import (
	repoPayment "payment-service/repositories/payment"
	repositories "payment-service/repositories/payment"
	repoHistory "payment-service/repositories/paymenthistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetPayment() repoPayment.IPaymentRepository
	GetPaymentHistory() repoHistory.IPaymentHistoryRepository
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetPayment() repositories.IPaymentRepository {
	return repoPayment.NewPaymentRepository(r.db)
}

func (r *Registry) GetPaymentHistory() repoHistory.IPaymentHistoryRepository {
	return repoHistory.NewPaymentHistoryRepository(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
