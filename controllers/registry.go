package controllers

import (
	fieldController "field-services/controllers/field"
	fieldScheduleController "field-services/controllers/fieldschedule"
	timeController "field-services/controllers/time"
	"field-services/services"
)

type Registry struct {
	service services.IServiceRegistry
}

// GetField implements [IControllerRegistry].
func (r *Registry) GetField() fieldController.IFieldController {
	return fieldController.NewFieldController(r.service)
}

// GetFieldSchedule implements [IControllerRegistry].
func (r *Registry) GetFieldSchedule() fieldScheduleController.IFieldScheduleController {
	return fieldScheduleController.NewFieldScheduleController(r.service)
}

// GetTime implements [IControllerRegistry].
func (r *Registry) GetTime() timeController.ITimeController {
	return timeController.NewTimeController(r.service)
}

type IControllerRegistry interface {
	GetField() fieldController.IFieldController
	GetFieldSchedule() fieldScheduleController.IFieldScheduleController
	GetTime() timeController.ITimeController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}
