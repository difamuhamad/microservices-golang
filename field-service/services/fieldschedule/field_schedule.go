package services

import (
	"context"
	"field-service/common/utils"
	"field-service/domain/dto"
	"field-service/repositories"
	"fmt"
	"time"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) *utils.PaginationResult
	GetAllByFieldIDAndDate(context.Context, *dto.FieldScheduleRequestParam, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleForBookingResponse, error)
	GeneratesScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}

// Get All Schedule with Pagination
func (f *FieldScheduleService) GetAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*utils.PaginationResult, error) {
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

// Convert month name to Indonesia
func (f *FieldScheduleService) converMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	IndonesiaMonth := map[string]string{
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
	formattedDate = fmt.Sprintf("%s %s", day, IndonesiaMonth[month])
	return formattedDate
}

func (f *FieldScheduleService) GetAllByFieldIDAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookingResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if error != nil {
		return nil, err
	}

	fieldSchedules, err := f.repository.GetField().FindAllByFieldIDAndDate()
}
