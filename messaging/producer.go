package messaging

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Producer struct {
	client   *sqs.Client
	queueURL string
}

func NewProducer(queueURL string) *Producer {
	return &Producer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}
}

func (p *Producer) Send(payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(messagingContext, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})

	return err
}
