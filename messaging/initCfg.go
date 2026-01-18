package messaging

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg aws.Config
var messagingContext context.Context

func InicializarMessaging(ctx context.Context) {
	messagingContext = ctx
	var err error
	cfg, err = config.LoadDefaultConfig(messagingContext)
	if err != nil {
		log.Fatal(err)
	}
}
