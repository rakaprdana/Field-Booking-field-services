package repositories

import (
	"context"
	"errors"
	errWrap "field-services/common/errors"
	constErr "field-services/constants/errors"
	constErrTime "field-services/constants/errors/time"
	"field-services/domains/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRegistry struct {
	db *gorm.DB
}

// Create implements [ITimeRegistry].
func (t *TimeRegistry) Create(ctx context.Context, time *models.Time) (*models.Time, error) {
	time.UUID = uuid.New()

	err := t.db.WithContext(ctx).Create(time).Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return time, nil
}

// FindAll implements [ITimeRegistry].
func (t *TimeRegistry) FindAll(ctx context.Context) ([]models.Time, error) {
	var times []models.Time

	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return times, nil
}

// FindByID implements [ITimeRegistry].
func (t *TimeRegistry) FindByID(ctx context.Context, id int) (*models.Time, error) {
	var time models.Time

	err := t.db.WithContext(ctx).Where("id = ?", id).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(constErrTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &time, nil
}

// FindByUUID implements [ITimeRegistry].
func (t *TimeRegistry) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	var time models.Time

	err := t.db.WithContext(ctx).Where("uuid = ?", uuid).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(constErrTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &time, nil
}

type ITimeRegistry interface {
	FindAll(context.Context) ([]models.Time, error)
	FindByUUID(context.Context, string) (*models.Time, error)
	FindByID(context.Context, int) (*models.Time, error)
	Create(context.Context, *models.Time) (*models.Time, error)
}

func NewTimeRegistry(db *gorm.DB) ITimeRegistry {
	return &TimeRegistry{db: db}
}
