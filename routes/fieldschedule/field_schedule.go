package routes

import (
	"field-services/clients"
	"field-services/constants"
	"field-services/controllers"
	"field-services/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

// Run implements [IFieldRoute].
func (f *FieldScheduleRoute) Run() {
	group := f.group.Group("/field/schedule")
	group.GET("", middlewares.AuthenticatedWithoutToken(), f.controller.GetFieldSchedule().GetAllByFieldIDAndDate)
	group.PATCH("/:uui", middlewares.AuthenticatedWithoutToken(), f.controller.GetFieldSchedule().UpdateStatus)

	group.Use(middlewares.Authenticated())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client), f.controller.GetFieldSchedule().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client), f.controller.GetFieldSchedule().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Create)
	group.POST("/one-month", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().GenerateScheduleForOneMonth)
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Update)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Delete)
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFielScheduleRoute(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}
