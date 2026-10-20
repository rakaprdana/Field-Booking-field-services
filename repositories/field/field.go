package repositories

import (
	"context"
	"errors"
	errWrap "field-services/common/errors"
	constErr "field-services/constants/errors"
	constErrField "field-services/constants/errors/field"
	"field-services/domains/dto"
	"field-services/domains/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRegistry struct {
	db *gorm.DB
}

// Create implements [IFieldRegistry].
func (f *FieldRegistry) Create(ctx context.Context, req *models.Field) (*models.Field, error) {
	field := models.Field{
		UUID:         uuid.New(),
		Code:         req.Code,
		Name:         req.Name,
		Images:       req.Images,
		PricePerHour: req.PricePerHour,
	}

	err := f.db.WithContext(ctx).Create(&field).Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &field, nil
}

// Delete implements [IFieldRegistry].
func (f *FieldRegistry) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Field{}).Error

	if err != nil {
		return errWrap.WrapError(constErr.ErrSQLError)
	}

	return nil
}

// FindAllWithPagination implements [IFieldRegistry].
func (f *FieldRegistry) FindAllWithPagination(
	ctx context.Context, param *dto.FieldRequestParam,
) ([]models.Field, int64, error) {

	var (
		fields []models.Field
		sort   string
		total  int64
	)

	order := "desc"
	if param.SortOrder != nil && *param.SortOrder != " " {
		order = *param.SortOrder
	}

	if param.SortColumn != nil && *param.SortColumn != "" {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, order)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fields).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(constErr.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Model(&models.Field{}).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(constErr.ErrSQLError)
	}
	return fields, total, nil
}

// FindByUUID implements [IFieldRegistry].
func (f *FieldRegistry) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var field models.Field

	err := f.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&field).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(constErrField.ErrFieldNotFound)
		}
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &field, nil
}

// FindWithoutPagination implements [IFieldRegistry].
func (f *FieldRegistry) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field

	err := f.db.
		WithContext(ctx).
		Find(&fields).Error

	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return fields, nil
}

// Update implements [IFieldRegistry].
func (f *FieldRegistry) Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error) {
	field := models.Field{
		Code:         req.Code,
		Name:         req.Name,
		Images:       req.Images,
		PricePerHour: req.PricePerHour,
	}

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&field).Error
	if err != nil {
		return nil, errWrap.WrapError(constErr.ErrSQLError)
	}

	return &field, nil
}

type IFieldRegistry interface {
	FindAllWithPagination(context.Context, *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.Field, error)
	FindByUUID(context.Context, string) (*models.Field, error)
	Create(context.Context, *models.Field) (*models.Field, error)
	Update(context.Context, string, *models.Field) (*models.Field, error)
	Delete(context.Context, string) error
}

func NewFieldRegistry(db *gorm.DB) IFieldRegistry {
	return &FieldRegistry{db: db}
}
