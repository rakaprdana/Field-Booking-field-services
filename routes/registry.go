package routes

import (
	"field-services/clients"
	"field-services/controllers"

	fieldRoute "field-services/routes/field"
	fieldScheduleRoute "field-services/routes/fieldschedule"
	timeRoute "field-services/routes/time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

// Serve implements [IRegistry].
func (r *Registry) Serve() {
	r.fieldRoute().Run()
	r.fieldScheduleRoute().Run()
	r.timeRoute().Run()
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) IRegistry {
	return &Registry{
		group:      group,
		controller: controller,
		client:     client,
	}
}

func (r *Registry) fieldRoute() fieldRoute.IFieldRoute {
	return fieldRoute.NewFieldRoute(r.group, r.controller, r.client)
}
func (r *Registry) fieldScheduleRoute() fieldScheduleRoute.IFieldScheduleRoute {
	return fieldScheduleRoute.NewFielScheduleRoute(r.group, r.controller, r.client)
}
func (r *Registry) timeRoute() timeRoute.ITimeRoute {
	return timeRoute.NewTimeRoute(r.group, r.controller, r.client)
}
