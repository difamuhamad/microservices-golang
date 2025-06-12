package dto

import "time"

type FieldScheduleRequest struct {
	FieldID string   `json:"fieldID" validate:"required"`
	Date    string   `json:"date" validate:"required"`
	TimeIDs []string `json:"timeIDs" validate:"required"`
}

type GenerateFieldScheduleForOneMonthRequest struct {
	FieldID string `json:"fieldID" validate:"required"`
}

type UpdateFieldScheduleRequest struct {
	Date   string `json:"date" validate:"required"`
	TimeID string `json:"timeID" validate:"required"`
}

type FieldScheduleResponse struct {
	UUID         string     `json:"uuid"`
	FieldName    string     `json:"fieldName"`
	PricePerHour int        `json:"pricePerHour"`
	Date         int        `json:"date"`
	Status       int        `json:"status"`
	Time         int        `json:"time"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type FieldScheduleForBookingResponse struct {
	UUID         string `json:"uuid"`
	PricePerHour int    `json:"pricePerHour"`
	Date         int    `json:"date"`
	Status       int    `json:"status"`
	Time         int    `json:"time"`
}

type FieldScheduleRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}

type FieldScheduleByFieldIDAndDateRequestParam struct {
	Date string `form:"date" validate:"require"`
}
