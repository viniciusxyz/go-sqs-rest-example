package consumers

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/viniciusxyz/go-sqs-rest-example/config"
	"github.com/viniciusxyz/go-sqs-rest-example/messaging"
	"github.com/viniciusxyz/go-sqs-rest-example/models"
)

func LogHandler(messages []string) error {
	var resp models.DepositoRequest

	for _, msg := range messages {
		err := json.Unmarshal([]byte(msg), &resp) //Deserialização do objeto recebido pelo SQS
		if err != nil {
			return err
		}
		slog.Info("Deposito processado com sucesso", "response", resp)
	}

	slog.Info("Total de mensagens processadas no lote: ", strconv.Itoa(len(messages)), "Lote finalizado")

	return nil
}

func InicializarConsumers() {
	if !config.ConsumersEnabled {
		slog.Warn("A env SQS_CONSUMERS_ENABLED está como false. Os consumers serão desativados")
		return
	}
	consumer := messaging.NewConsumer(config.QueueURL)
	go consumer.StartConsumer(LogHandler) //Cria uma coroutine a parte para execução do consumer para não travar a aplicação
}
