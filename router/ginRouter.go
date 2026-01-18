package router

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciusxyz/go-sqs-rest-example/controllers"
)

func InicializarRotas() error {
	r := gin.Default()

	r.POST("/deposito", controllers.DepositoController)

	return r.Run(":8080")
}
