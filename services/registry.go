package services

import (
	"field-services/common/gcs"
	"field-services/repositories"
	fieldService "field-services/services/field"
	fieldScheduleService "field-services/services/fieldschedule"
	timeService "field-services/services/time"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

// GetField implements [IServiceRegistry].
func (r *Registry) GetField() fieldService.IFieldService {
	return fieldService.NewFieldService(r.repository, r.gcs)
}

// GetFieldSchedule implements [IServiceRegistry].
func (r *Registry) GetFieldSchedule() fieldScheduleService.IFieldScheduleService {
	return fieldScheduleService.NewFieldScheduleRegistry(r.repository)
}

// GetTime implements [IServiceRegistry].
func (r *Registry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(r.repository)
}

type IServiceRegistry interface {
	GetField() fieldService.IFieldService
	GetFieldSchedule() fieldScheduleService.IFieldScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &Registry{
		repository: repository,
		gcs:        gcs,
	}
}
