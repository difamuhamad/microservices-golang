package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	"field-service/constants"
	errConstant "field-service/constants/error"
	errField "field-service/constants/error/field"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

type IFieldScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, *[]models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, constants.FieldScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

// Find All Field with Pagination
func (f *FieldScheduleRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fieldsSchedule []models.FieldSchedule
		total          int64
	)

	// Default sorting
	sort := "created_at desc"
	if param.SortColumn != nil && param.SortOrder != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	}

	// Get total count first
	if err := f.db.WithContext(ctx).Model(&models.FieldSchedule{}).Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Get paginated results
	if err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Limit(param.Limit).
		Offset((param.Page - 1) * param.Limit).
		Order(sort).
		Find(&fieldsSchedule).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldsSchedule, total, nil
}

// Find All Field by ID and Date
func (f *FieldScheduleRepository) FindAllByFieldIDAndDate(ctx context.Context, fieldID int, date string) ([]models.FieldSchedule, error) {
	var fieldSchedules []models.FieldSchedule
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", fieldID).
		Where("date = ?", date).
		Joins("LEFT JOIN times ON field_schedules.time_id = times.id").
		Order("times.start_time asc").
		Find(&fieldSchedules).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return fieldSchedules, nil
}

// Find Field by UUID
func (f *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil
}

// Find Field by Date and Time
func (f *FieldScheduleRepository) FindByDateAndTimeID(ctx context.Context, date string, timeID int, fieldID int) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("date = ?", date).
		Where("time_id = ?", timeID).
		Where("field_id = ?", fieldID).
		First(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil
}

// Create Field Schedules
func (f *FieldScheduleRepository) Create(ctx context.Context, req *[]models.FieldSchedule) error {
	if err := f.db.WithContext(ctx).Create(req).Error; err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}

// Update Field
func (f *FieldScheduleRepository) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Update all fields from req
	fieldSchedule.Date = req.Date
	fieldSchedule.TimeID = req.TimeID
	fieldSchedule.FieldID = req.FieldID
	fieldSchedule.Status = req.Status
	// Add other fields as needed

	if err := f.db.WithContext(ctx).Save(&fieldSchedule).Error; err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldSchedule, nil
}

// Update Field Status
func (f *FieldScheduleRepository) UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error {
	if err := f.db.WithContext(ctx).
		Model(&models.FieldSchedule{}).
		Where("uuid = ?", uuid).
		Update("status", status).Error; err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}

// Delete Field
func (f *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	var existing models.FieldSchedule
	if err := f.db.WithContext(ctx).Where("uuid = ?", uuid).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errField.ErrFieldScheduleNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	if err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}).Error; err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}
