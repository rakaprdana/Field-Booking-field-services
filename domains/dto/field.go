package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Name         string                 `form:"name" validate:"required"`
	Code         string                 `form:"code" validate:"required"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images"` //Unable validation required
}

type UpdateFieldRequest struct {
	UUID         string                 `form:"uuid" validate:"required"`
	Name         string                 `form:"name" validate:"required"`
	Code         string                 `form:"code" validate:"required"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images"`
}

type FieldResponse struct {
	UUID         uuid.UUID  `json:"uuid"`
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	PricePerHour int        `json:"pricePerHour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type FieldDetailResponse struct {
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	PricePerHour int        `json:"pricePerHour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

// Pagination
type FieldRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}
