package services

import (
	"br.com.viniciusxyz/go-sqs-rest-example/config"
	"br.com.viniciusxyz/go-sqs-rest-example/messaging"
	"br.com.viniciusxyz/go-sqs-rest-example/models"
)

var producer *messaging.Producer

func SendDeposito(deposito models.DepositoRequest) error {

	if producer == nil {
		producer = messaging.NewProducer(config.QueueURL)
	}

	err := producer.Send(deposito)
	return err
}
