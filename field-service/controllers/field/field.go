package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldController struct {
	service services.IServiceRegistry
}

type IFieldController interface {
	GetAllWithPagination(*gin.Context)
	GetAllWithoutPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewFieldController(service services.IServiceRegistry) IFieldController {
	return &FieldController{service: service}
}

// Get All With Pagination Controller
func (f *FieldController) GetAllWithPagination(c *gin.Context) {
	var params dto.FieldRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetField().GetAllWithPagination(c, &params)
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

// Get All without Pagination Controller
func (f *FieldController) GetAllWithoutPagination(c *gin.Context) {

	result, err := f.service.GetField().GetAllWithoutPagination(c)
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

// Get All By UUID Controller
func (f *FieldController) GetByUUID(c *gin.Context) {

	result, err := f.service.GetField().GetByUUID(c, c.Param("uuid"))
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

// Create Field Controller
func (f *FieldController) Create(c *gin.Context) {
	var request dto.FieldRequest
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetField().Create(c, &request)
	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}

// Update Field Controller
func (f *FieldController) Update(c *gin.Context) {
	var request dto.UpdateFieldRequest
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetField().Update(c, c.Param("uuid"), &request)
	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}

// Delete Field Controller
func (f *FieldController) Delete(c *gin.Context) {
	err := f.service.GetField().Delete(c, c.Param("uuid"))
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
		Gin:  c,
	})

}
