package config

import (
	"os"
	"strconv"
)

var QueueURL = os.Getenv("SQS_QUEUE_URL")
var ConsumersEnabled = getEnvBool("SQS_CONSUMERS_ENABLED", true)

// Função para buscar uma variável de ambiente definindo o valor padrão de true
func getEnvBool(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return b
}
