package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldScheduleController struct {
	service services.IServiceRegistry
}

type IFieldScheduleController interface {
	GetAllWithPagination(*gin.Context)
	GetAllByFieldIDAndDate(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	UpdateStatus(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GenerateScheduleForOneMonth(*gin.Context)
}

func NewScheduleController(service services.IServiceRegistry) IFieldScheduleController {
	return &FieldScheduleController{service: service}
}

// Get All Schedule With Pagination Controller
func (f *FieldScheduleController) GetAllWithPagination(c *gin.Context) {
	var params dto.FieldScheduleRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// Validate the request first
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

	result, err := f.service.GetFieldSchedule().GetAllWithPagination(c, &params)
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

// Get All Schedule by Field ID and Date
func (f *FieldScheduleController) GetAllByFieldIDAndDate(c *gin.Context) {
	var params dto.FieldScheduleByFieldIDAndDateRequestParam
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

	result, err := f.service.GetFieldSchedule().GetAllByFieldIDAndDate(c, c.Param("uuid"), params.Date)
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

// Get Schedule by UUID Controller
func (f *FieldScheduleController) GetByUUID(c *gin.Context) {
	result, err := f.service.GetFieldSchedule().GetByUUID(c, c.Param("uuid"))
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

// Create Field Schedule Controller
func (f *FieldScheduleController) Create(c *gin.Context) {
	var params dto.FieldScheduleRequest
	err := c.ShouldBindJSON(&params)
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

	err = f.service.GetFieldSchedule().Create(c, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  c,
	})

}

// Update Schedule Controller
func (f *FieldScheduleController) Update(c *gin.Context) {
	var params dto.UpdateFieldScheduleRequest
	err := c.ShouldBindJSON(&params)
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

	result, err := f.service.GetFieldSchedule().Update(c, c.Param("uuid"), &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Data: result,
		Gin:  c,
	})

}

// Update Schedule Status Controller
func (f *FieldScheduleController) UpdateStatus(c *gin.Context) {
	var request dto.UpdateStatusFieldScheduleRequest
	err := c.ShouldBindJSON(&request)
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

	err = f.service.GetFieldSchedule().UpdateStatus(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  c,
	})

}

// Delete Schedule Controller
func (f *FieldScheduleController) Delete(c *gin.Context) {
	err := f.service.GetFieldSchedule().Delete(c, c.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
	}
	return
}

// Generate Schedule for one month Controller
func (f *FieldScheduleController) GenerateScheduleForOneMonth(c *gin.Context) {
	var params dto.GenerateFieldScheduleForOneMonthRequest
	err := c.ShouldBindJSON(&params)
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

	err = f.service.GetFieldSchedule().GenerateScheduleForOneMonth(c, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  c,
	})

}
