package router

import (
	"br.com.viniciusxyz/go-sqs-rest-example/controllers"
	"github.com/gin-gonic/gin"
)

func InicializarRotas() error {
	r := gin.Default()

	r.POST("/deposito", controllers.DepositoController)

	return r.Run(":8080")
}
