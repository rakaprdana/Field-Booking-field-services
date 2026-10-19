package repositories

import (
	"context"
	"errors"
	errWrap "field-services/common/errors"
	"field-services/constants"
	constErr "field-services/constants/errors"
	constErrFieldSchedule "field-services/constants/errors/fieldschedule"
	"field-services/domains/dto"
	"field-services/domains/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRegistry struct {
	db *gorm.DB
}

// Create implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) Create(ctx context.Context, req []models.FieldSchedule) error {
	err := f.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

// Delete implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}).Error

	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

// FindAllWithPagination implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) FindAllWithPagination(
	ctx context.Context, param *dto.FieldScheduleRequestParam,
) ([]models.FieldSchedule, int64, error) {

	var (
		fieldSchedules []models.FieldSchedule
		sort           string
		total          int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fieldSchedules).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(constErr.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Order(sort).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(constErr.ErrSQLError)
	}

	return fieldSchedules, total, nil
}

// FindByUUID implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedules models.FieldSchedule

	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&fieldSchedules).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(constErrFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &fieldSchedules, nil
}

// FindByDateAndTimeID implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) FindByDateAndTimeID(
	ctx context.Context,
	date string,
	timeID int,
	fieldID int) (*models.FieldSchedule, error) {
	var fieldSchedules models.FieldSchedule

	err := f.db.
		WithContext(ctx).
		Where("date = ?", date).
		Where("field_id = ?", fieldID).
		Where("time_id = ?", timeID).
		First(&fieldSchedules).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(constErrFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &fieldSchedules, nil
}

// FindAllByFieldIDAndDate implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) FindAllByFieldIDAndDate(
	ctx context.Context,
	fieldID int,
	date string) ([]models.FieldSchedule, error) {
	var fieldSchedules []models.FieldSchedule

	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", fieldID).
		Where("date = ?", date).
		Joins("LEFT JOIN times ON field_schedules.time_id = times.id").
		Order("times.start_time asc").
		Find(&fieldSchedules).
		Error

	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return fieldSchedules, nil
}

// Update implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) Update(
	ctx context.Context,
	uuid string,
	req *models.FieldSchedule,
) (*models.FieldSchedule, error) {

	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return fieldSchedule, nil
}

// UpdateStatus implements [IFieldScheduleRegistry].
func (f *FieldScheduleRegistry) UpdateStatus(
	ctx context.Context,
	uuid string,
	status constants.FieldScheduleStatus,
) error {

	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	fieldSchedule.Status = status
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

type IFieldScheduleRegistry interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, []models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, string, constants.FieldScheduleStatus) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRegistry(db *gorm.DB) IFieldScheduleRegistry {
	return &FieldScheduleRegistry{db: db}
}
