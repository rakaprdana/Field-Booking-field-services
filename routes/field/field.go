package routes

import (
	"field-services/clients"
	"field-services/constants"
	"field-services/controllers"
	"field-services/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

// Run implements [IFieldRoute].
func (f *FieldRoute) Run() {
	group := f.group.Group("/field")
	group.GET("", middlewares.AuthenticatedWithoutToken(), f.controller.GetField().GetAllWithoutPagination)
	group.GET("/:uuid", middlewares.AuthenticatedWithoutToken(), f.controller.GetField().GetByUUID)

	group.Use(middlewares.Authenticated())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client), f.controller.GetField().GetAllWithPagination)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetField().Create)
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetField().Update)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetField().Delete)
}

type IFieldRoute interface {
	Run()
}

func NewFieldRoute(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry) IFieldRoute {
	return &FieldRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}
