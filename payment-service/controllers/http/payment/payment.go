package controllers

import (
	"net/http"
	errValidation "payment-service/common/error"
	"payment-service/common/response"
	"payment-service/domain/dto"
	"payment-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PaymentController struct {
	service services.IServiceRegistry
}

type IPaymentController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Webhook(*gin.Context)
}

func NewPaymentController(service services.IServiceRegistry) IPaymentController {
	return &PaymentController{service: service}
}

func (p *PaymentController) GetAllWithPagination(c *gin.Context) {
	var param dto.PaymentRequestParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// Validate param request
	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := p.service.GetPayment().GetAllWithPagination(c, &param)
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

func (p *PaymentController) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	result, err := p.service.GetPayment().GetByUUID(c, uuid)
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

func (p *PaymentController) Create(c *gin.Context) {
	var request dto.PaymentRequest

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
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := p.service.GetPayment().Create(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated, //use status created
		Data: result,
		Gin:  c,
	})
}

func (p *PaymentController) Webhook(c *gin.Context) {
	var request dto.Webhook

	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	err = p.service.GetPayment().Webhook(c, &request)
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
