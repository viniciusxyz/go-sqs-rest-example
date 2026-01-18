package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/viniciusxyz/go-sqs-rest-example/consumers"
	"github.com/viniciusxyz/go-sqs-rest-example/messaging"
	"github.com/viniciusxyz/go-sqs-rest-example/router"
)

func main() {

	//Cria o contexto para gracefulShutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	messaging.InicializarMessaging(ctx) //Inicializa as configurações da AWS com o context atual
	consumers.InicializarConsumers()    //Inicializa os consumers

	go func() { //Executa o Gin em uma coroutine separada
		if err := router.InicializarRotas(); err != nil {
			slog.Error("Servidor Gin desligado", "error", err)
			stop() //Quando o servidor do GIN termina ai o stop() é chamado
		}
	}()

	<-ctx.Done() //Aguarda até o cancelamento do contexto
	slog.Info("Aplicação desligada...")
}
