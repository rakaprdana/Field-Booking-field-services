package models

import (
	"field-services/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldSchedule struct {
	ID        uint                          `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID                     `gorm:"type:uuid;not null"`
	FieldID   uint                          `gorm:"type:int;not null"`
	TimeID    uint                          `gorm:"type:int;not null"`
	Date      time.Time                     `gorm:"type:int;not null"`
	Status    constants.FieldScheduleStatus `gorm:"type:int;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt

	Field Field `gorm:"foreignKey:field_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Time  Time  `gorm:"foreignKey:time_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
