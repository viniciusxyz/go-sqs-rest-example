package messaging

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
	workers  int
}

func NewConsumer(queueURL string, workers int) *Consumer {
	return &Consumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
		workers:  workers,
	}
}

type HandlerFunc func(message types.Message, workerId int) error

func (c *Consumer) worker(
	ctx context.Context,
	id int,
	jobs <-chan types.Message,
	handler HandlerFunc,
) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("worker finalizado", "id", id)
			return

		case msg, ok := <-jobs:
			if !ok {
				return
			}

			if err := handler(msg, id); err != nil {
				slog.Error("erro ao processar mensagem", "error", err)
				continue
			}

			// delete somente após sucesso
			_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				slog.Error("erro ao deletar mensagem", "error", err)
			}
		}
	}
}

func (c *Consumer) receiveLoop(
	ctx context.Context,
	jobs chan<- types.Message,
) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("receiver finalizado")
			return
		default:
		}

		resp, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 10, // Quantidade máxima de mensagens recebidas
			WaitTimeSeconds:     20, // Define um tempo de espera para não ter que consultar o SQS em loop constante
			VisibilityTimeout:   30, // Tempo máximo para o processamento da mensagem
		})
		if err != nil {
			slog.Error("erro ao receber mensagens", "error", err)
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range resp.Messages {
			select {
			case jobs <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Inicia o consumo executando o HandlerFunc informado
func (c *Consumer) StartConsumer(handler HandlerFunc) error {

	jobs := make(chan types.Message)

	var wg sync.WaitGroup

	// workers
	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.worker(messagingContext, id, jobs, handler)
		}(i)
	}

	// receiver
	go func() {
		defer close(jobs)
		c.receiveLoop(messagingContext, jobs)
	}()

	<-messagingContext.Done()
	slog.Info("consumer finalizando, aguardando workers...")
	wg.Wait()
	return nil
}
