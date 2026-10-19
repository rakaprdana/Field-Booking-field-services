package services

import (
	"bytes"
	"context"
	"field-services/common/gcs"
	"field-services/common/utils"
	constErr "field-services/constants/errors"
	"field-services/domains/dto"
	"field-services/domains/models"
	"field-services/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"
)

type FieldService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

// Create implements [IFieldService].
func (f *FieldService) Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error) {
	imageUrl, err := f.uploadImage(ctx, req.Images)
	if err != nil {
		return nil, err
	}

	field, err := f.repository.GetField().Create(ctx, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	return response, nil

}

// Delete implements [IFieldService].
func (f *FieldService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetField().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

// GetAllWithPagination implements [IFieldService].
func (f *FieldService) GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*utils.PaginationResult, error) {
	fields, total, err := f.repository.GetField().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResult = append(fieldResult, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	pagination := &utils.PaginationParams{
		Count: int(total),
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldResult,
	}
	response := utils.GeneratePagination(*pagination)
	return &response, nil
}

// GetAllWithoutPagination implements [IFieldService].
func (f *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResult = append(fieldResult, dto.FieldResponse{
			UUID:         field.UUID,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
		})
	}

	return fieldResult, nil
}

// GetByUUID implements [IFieldService].
func (f *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldResult := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	return &fieldResult, nil
}

// Update implements [IFieldService].
func (f *FieldService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var imageUrl []string

	if req.Images == nil {
		imageUrl = field.Images
	} else {
		imageUrl, err = f.uploadImage(ctx, req.Images)
		if err != nil {
			return nil, err
		}
	}

	fieldResult, err := f.repository.GetField().Update(ctx, uuid, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.FieldResponse{
		UUID:         fieldResult.UUID,
		Code:         fieldResult.Code,
		Name:         fieldResult.Name,
		PricePerHour: fieldResult.PricePerHour,
		Images:       fieldResult.Images,
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}

	return response, nil
}

type IFieldService interface {
	GetAllWithPagination(context.Context, *dto.FieldRequestParam) (*utils.PaginationResult, error)
	GetAllWithoutPagination(context.Context) ([]dto.FieldResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldResponse, error)
	Create(context.Context, *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(context.Context, string, *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(context.Context, string) error
}

func NewFieldService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IFieldService {
	return &FieldService{
		repository: repository,
		gcs:        gcs,
	}
}

func (f *FieldService) validateUpload(images []multipart.FileHeader) error {
	if images == nil || len(images) == 0 {
		return constErr.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return constErr.ErrSizeTooBig
		}
	}
	return nil
}

func (f *FieldService) processAndUploadImage(
	ctx context.Context,
	image multipart.FileHeader,
) (string, error) {
	file, err := image.Open()

	if err != nil {
		return "", nil
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("images/%s-%s-%s",
		time.Now().Format("20261019131502"),
		image.Filename,
		path.Ext(image.Filename))

	url, err := f.gcs.UploadFile(ctx, filename, buffer.Bytes())
	if err != nil {
		return "", nil
	}

	return url, nil
}

func (f *FieldService) uploadImage(
	ctx context.Context,
	images []multipart.FileHeader,
) ([]string, error) {
	err := f.validateUpload(images)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := f.processAndUploadImage(ctx, image)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}
	return urls, nil
}
