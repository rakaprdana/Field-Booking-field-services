package controllers

import (
	errValidate "field-services/common/errors"
	response "field-services/common/reponses"
	"field-services/domains/dto"
	"net/http"

	"field-services/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TimeController struct {
	service services.IServiceRegistry
}

// Create implements [ITimeController].
func (t *TimeController) Create(ctx *gin.Context) {
	var req dto.TimeRequest
	err := ctx.ShouldBindJSON(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidate.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}
	result, err := t.service.GetTime().Create(ctx, &req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// GetAll implements [ITimeController].
func (t *TimeController) GetAll(ctx *gin.Context) {
	result, err := t.service.GetTime().GetAll(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// GetByUUID implements [ITimeController].
func (t *TimeController) GetByUUID(ctx *gin.Context) {
	result, err := t.service.GetTime().GetByUUID(ctx, ctx.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

type ITimeController interface {
	GetAll(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
}

func NewTimeController(service services.IServiceRegistry) ITimeController {
	return &TimeController{service: service}
}
