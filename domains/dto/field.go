package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieladRequest struct {
	Name         string                 `json:"name" validate:"required"`
	Code         string                 `json:"code" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images"` //Unable validation required
}

type UpdateFieldRequest struct {
	UUID         string                 `json:"uuid" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	Code         string                 `json:"code" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images"`
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

type FieladDetailResponse struct {
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
