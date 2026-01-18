package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciusxyz/go-sqs-rest-example/models"
	"github.com/viniciusxyz/go-sqs-rest-example/services"
)

func DepositoController(c *gin.Context) {

	var req models.DepositoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Payload inválido",
		})
		return
	}

	if err := services.SendDeposito(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to send message to SQS",
		})
		return
	}

	c.JSON(http.StatusOK, models.DepositoResponse{
		Status:   "OK",
		Mensagem: "Depósito realizado com sucesso",
	})
}
