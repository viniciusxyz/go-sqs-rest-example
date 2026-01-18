package models

type DepositoRequest struct {
	ContaOrigem  string  `json:"contaOrigem" example:"100000"`
	ContaDestino string  `json:"contaDestino" example:"100001"`
	Valor        float64 `json:"valor" example:"1000"`
}

type DepositoResponse struct {
	Status   string `json:"status" example:"OK"`
	Mensagem string `json:"mensagem" example:"Depósito realizado com sucesso"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Payload inválido"`
}
