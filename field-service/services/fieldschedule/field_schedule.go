package services

import (
	"context"
	"field-service/common/utils"
	"field-service/constants"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*utils.PaginationResult, error)
	GetAllByFieldIDAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}

// Get All Field with Pagination
func (f *FieldScheduleService) GetAllWithPagination(
	ctx context.Context,
	param *dto.FieldScheduleRequestParam,
) (*utils.PaginationResult, error) {
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldScheduleResults := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldScheduleResults = append(fieldScheduleResults, dto.FieldScheduleResponse{
			UUID:         schedule.UUID,
			FieldName:    schedule.Field.Name,
			Date:         schedule.Date.Format("2006-01-02"),
			PricePerHour: schedule.Field.PricePerHour,
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:    schedule.CreatedAt,
			UpdatedAt:    schedule.UpdatedAt,
		})
	}

	pagination := &utils.PaginationParam{
		Count: total,
		Limit: param.Limit,
		Page:  param.Page,
		Data:  fieldScheduleResults,
	}

	response := utils.GeneratePagination(*pagination)
	return &response, nil
}

func (f *FieldScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	indonesiaMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Agu",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formattedDate := date.Format("02 Jan")
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indonesiaMonth[month])
	return formattedDate
}

// Get All Field by ID and Date
func (f *FieldScheduleService) GetAllByFieldIDAndDate(
	ctx context.Context,
	uuid, date string,
) ([]dto.FieldScheduleForBookingResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	fieldScheduleResults := make([]dto.FieldScheduleForBookingResponse, 0, len(fieldSchedules))
	for _, fieldSchedule := range fieldSchedules {
		pricePerHour := float64(fieldSchedule.Field.PricePerHour)
		startTime, _ := time.Parse("15:04:05", fieldSchedule.Time.StartTime)
		endTime, _ := time.Parse("15:04:05", fieldSchedule.Time.EndTime)
		fieldScheduleResults = append(fieldScheduleResults, dto.FieldScheduleForBookingResponse{
			UUID:         fieldSchedule.UUID,
			PricePerHour: utils.RupiahFormat(&pricePerHour),
			Date:         f.convertMonthName(fieldSchedule.Date.Format("2006-01-02")),
			Status:       fieldSchedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04")),
		})
	}

	return fieldScheduleResults, nil
}

// Get All Field by UUID
func (f *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Date:         fieldSchedule.Date.Format(time.DateOnly),
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}
	return &response, nil
}

// Create Field Data for one month
func (f *FieldScheduleService) GenerateScheduleForOneMonth(
	ctx context.Context,
	request *dto.GenerateFieldScheduleForOneMonthRequest,
) error {
	field, err := f.repository.GetField().FindByUUID(ctx, request.FieldID)
	if err != nil {
		return err
	}

	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDays := 30
	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays)
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)
	for i := 0; i < numberOfDays; i++ {
		currentDate := now.AddDate(0, 0, i)
		for _, item := range times {
			schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
				ctx,
				currentDate.Format(time.DateOnly),
				int(item.ID),
				int(field.ID),
			)
			if err != nil {
				return err
			}

			if schedule != nil {
				return errFieldSchedule.ErrFieldScheduleIsExist
			}

			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    currentDate,
				Status:  constants.Available,
			})
		}
	}

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil
}

// Create Field Data
func (f *FieldScheduleService) Create(ctx context.Context, request *dto.FieldScheduleRequest) error {
	field, err := f.repository.GetField().FindByUUID(ctx, request.FieldID)
	if err != nil {
		return err
	}

	fieldSchedules := make([]models.FieldSchedule, 0, len(request.TimeIDs))
	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	for _, timeID := range request.TimeIDs {
		scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		// Check if the field already exist
		schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			return err
		}

		if schedule != nil {
			return errFieldSchedule.ErrFieldScheduleIsExist
		}

		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParsed,
			Status:  constants.Available,
		})
	}

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil
}

// Update Field Data
func (f *FieldScheduleService) Update(
	ctx context.Context,
	uuid string,
	request *dto.UpdateFieldScheduleRequest,
) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, request.TimeID)
	if err != nil {
		return nil, err
	}

	isTimeExist, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
		ctx,
		request.Date,
		int(scheduleTime.ID),
		int(fieldSchedule.FieldID),
	)
	if err != nil {
		return nil, err
	}

	if isTimeExist != nil && request.Date != fieldSchedule.Date.Format(time.DateOnly) {
		checkDate, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
			ctx,
			request.Date,
			int(scheduleTime.ID),
			int(fieldSchedule.FieldID),
		)
		if err != nil {
			return nil, err
		}

		if checkDate != nil {
			return nil, errFieldSchedule.ErrFieldScheduleIsExist
		}
	}

	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	fieldResult, err := f.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldResult.UUID,
		FieldName:    fieldResult.Field.Name,
		Date:         fieldResult.Date.Format(time.DateOnly),
		PricePerHour: fieldResult.Field.PricePerHour,
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}
	return &response, nil
}

// Update Field Status
func (f *FieldScheduleService) UpdateStatus(
	ctx context.Context,
	request *dto.UpdateStatusFieldScheduleRequest,
) error {
	for _, item := range request.FieldScheduleIDs {
		_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = f.repository.GetFieldSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Field Data
func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
