package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFieldScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldScheduleRoute) Run() {
	// Schedule Routes Group
	group := f.group.Group("/field/schedule")

	// Without login routes :
	group.GET("", middlewares.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().GetAllByFieldIDAndDate)

	group.GET("/:uuid", middlewares.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().GetByUUID)

	group.PATCH("", middlewares.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().UpdateStatus)

	// Must login routes :
	group.Use(middlewares.Authenticate())

	group.GET("/pagination", middlewares.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetFieldSchedule().GetAllWithPagination)

	group.POST("", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Create)

	group.POST("/one-month", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().GenerateScheduleForOneMonth)

	group.PUT("/:uuid", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Update)

	group.DELETE("/:uuid", middlewares.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Delete)
}
