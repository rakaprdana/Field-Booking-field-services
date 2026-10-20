package controllers

import (
	response "field-services/common/reponses"
	"field-services/domains/dto"
	"net/http"

	"field-services/services"

	"github.com/gin-gonic/gin"
)

type TimeController struct {
	service services.IServiceRegistry
}

// Create implements [ITimeController].
func (t *TimeController) Create(ctx *gin.Context) {
	var req dto.TimeRequest
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
