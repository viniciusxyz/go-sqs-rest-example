package consumers

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/viniciusxyz/go-sqs-rest-example/config"
	"github.com/viniciusxyz/go-sqs-rest-example/messaging"
	"github.com/viniciusxyz/go-sqs-rest-example/models"
)

func LogHandler(message types.Message, workerId int) error {
	var resp models.DepositoRequest

	err := json.Unmarshal([]byte(*message.Body), &resp) //Deserialização do objeto recebido pelo SQS
	if err != nil {
		return err
	}

	slog.Info("worker", strconv.Itoa(workerId), "Deposito processado com sucesso", "response", resp)
	return nil
}

func InicializarConsumers() {
	if !config.ConsumersEnabled {
		slog.Warn("A env SQS_CONSUMERS_ENABLED está como false. Os consumers serão desativados")
		return
	}
	consumer := messaging.NewConsumer(config.QueueURL, 10)
	go consumer.StartConsumer(LogHandler) //Cria uma coroutine a parte para execução do consumer para não travar a aplicação
}
