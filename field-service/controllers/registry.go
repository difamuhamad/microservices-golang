package controllers

import (
	fieldController "field-service/controllers/field"
	fieldScheduleController "field-service/controllers/fieldschedule"
	timeController "field-service/controllers/time"
	"field-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetField() fieldController.IFieldController
	GetFieldSchedule() fieldScheduleController.IFieldScheduleController
	GetTime() timeController.ITimeController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetField() fieldController.IFieldController {
	return fieldController.NewFieldController(r.service)
}

func (r *Registry) GetFieldSchedule() fieldScheduleController.IFieldScheduleController {
	return fieldScheduleController.NewScheduleController(r.service)
}

func (r *Registry) GetTime() timeController.ITimeController {
	return timeController.NewFieldController(r.service)
}
