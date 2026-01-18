package messaging

import (
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

type HandlerFunc func(body string) error

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

		resp, err := c.client.ReceiveMessage(messagingContext, &sqs.ReceiveMessageInput{
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

		if len(resp.Messages) == 0 {
			continue
		}

		for _, msg := range resp.Messages {
			if msg.Body == nil {
				slog.Warn("Mensagem sem body recebida")
				continue
			}

			// Processa a mensagem
			if err := handler(*msg.Body); err != nil {
				slog.Error("Erro ao processar mensagem", "error", err)
				continue // NÃO deleta → SQS reentrega
			}

			// Deleta somente se processar com sucesso
			_, err := c.client.DeleteMessage(messagingContext, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(c.queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				slog.Error("Erro ao deletar mensagem do SQS", "error", err)
			}
		}
	}
}
