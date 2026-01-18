package config

import "os"

var QueueURL = os.Getenv("SQS_QUEUE_URL")
