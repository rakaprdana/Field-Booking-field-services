package routes

import (
	"field-services/clients"
	"field-services/constants"
	"field-services/controllers"
	"field-services/middlewares"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

// Run implements [ITimeRoute].
func (t *TimeRoute) Run() {
	group := t.group.Group("/time")
	group.Use(middlewares.Authenticated())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client), t.controller.GetTime().GetAll)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client), t.controller.GetTime().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client), t.controller.GetTime().Create)
}

type ITimeRoute interface {
	Run()
}

func NewTimeRoute(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) ITimeRoute {
	return &TimeRoute{
		controller: controller,
		client:     client,
		group:      group,
	}
}
