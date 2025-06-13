package controllers

import (
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TimeController struct {
	service services.IServiceRegistry
}

type ITimeController interface {
	GetAll(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
}

func NewFieldController(service services.IServiceRegistry) ITimeController {
	return &TimeController{service: service}
}

// Get All Time
func (t *TimeController) GetAll(c *gin.Context) {
	result, err := t.service.GetTime().GetAll(c)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}

// Get Time by UUID
func (t *TimeController) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	result, err := t.service.GetTime().GetByUUID(c, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}

// Create Time Controller
func (t *TimeController) Create(c *gin.Context) {
	var request dto.TimeRequest
	result, err := t.service.GetTime().Create(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}
