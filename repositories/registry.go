package repositories

import (
	repoFieldSchedule "field-services/repositories/fieladschedule"
	repoField "field-services/repositories/field"
	repoTime "field-services/repositories/time"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

// GetField implements [IRepositoryRegistry].
func (r *Registry) GetField() repoField.IFieldRegistry {
	return repoField.NewFieldRegistry(r.db)
}

// GetFieldSchedule implements [IRepositoryRegistry].
func (r *Registry) GetFieldSchedule() repoFieldSchedule.IFieldScheduleRegistry {
	return repoFieldSchedule.NewFieldScheduleRegistry(r.db)
}

// GetTime implements [IRepositoryRegistry].
func (r *Registry) GetTime() repoTime.ITimeRegistry {
	return repoTime.NewTimeRegistry(r.db)
}

type IRepositoryRegistry interface {
	GetField() repoField.IFieldRegistry
	GetFieldSchedule() repoFieldSchedule.IFieldScheduleRegistry
	GetTime() repoTime.ITimeRegistry
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}
