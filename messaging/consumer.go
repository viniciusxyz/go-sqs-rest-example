package messaging

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
}

func NewConsumer(queueURL string) *Consumer {
	return &Consumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}
}

type HandlerFunc func(messages []string) error

// Inicia o consumo executando o HandlerFunc informado
func (c *Consumer) StartConsumer(handler HandlerFunc) error {
	slog.Info("SQS consumer iniciado", "queue", c.queueURL)

	for {
		// Valida se o Contexto foi finalizado para parar o consumo de mensagens
		select {
		case <-messagingContext.Done():
			slog.Info("SQS consumer finalizado")
			return messagingContext.Err()
		default:
		}

		out, err := c.client.ReceiveMessage(messagingContext, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL), //URL da file do SQS
			MaxNumberOfMessages: 10,                     // Quantidade máxima de mensagens recebidas
			WaitTimeSeconds:     20,                     // Define um tempo de espera para não ter que consultar o SQS em loop constante
			VisibilityTimeout:   30,                     // Tempo máximo para o processamento da mensagem
		})

		if err != nil {
			slog.Error("Erro ao receber mensagens do SQS", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(out.Messages) == 0 {
			continue
		}

		bodies := make([]string, 0, len(out.Messages)) // Cria um array que comportará os bodies
		for _, msg := range out.Messages {
			bodies = append(bodies, aws.ToString(msg.Body)) // Adiciona as mensagens no array as transformando em strings
		}

		if err := handler(bodies); err != nil {
			slog.Error("erro ao processar lote", "error", err)
			continue
		}

		c.deleteBatch(out.Messages)
	}
}

func (c *Consumer) deleteBatch(
	messages []types.Message,
) {

	entries := make([]types.DeleteMessageBatchRequestEntry, 0, len(messages)) // Cria um array de Entrys (mensagens)

	for i, msg := range messages {
		entries = append(entries, types.DeleteMessageBatchRequestEntry{ // Popula o array
			Id:            aws.String(fmt.Sprintf("msg-%d", i)),
			ReceiptHandle: msg.ReceiptHandle,
		})
	}

	_, err := c.client.DeleteMessageBatch(messagingContext, &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.queueURL),
		Entries:  entries, // Passa o array inteiro para delete
	})

	if err != nil {
		slog.Error("erro ao deletar batch", "error", err)
	}
}
