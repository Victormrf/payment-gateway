package domain

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found") // retornado quando uma conta não é encontrada
	ErrDuplicatedAPIKey = errors.New("api key already exists") // retornado em tentativa de criar quando uma chave de API já existe
	ErrInvoiceNotFound = errors.New("invoice not found") // retornado quando uma fatura não é encontrada
	ErrUnauthorizedAccess = errors.New("unauthorized access") // retornado quando o acesso não é autorizado
	ErrInvalidAmount = errors.New("amount must be greater than 0") // retornado quando o valor da fatura é inválido
	ErrInvalidStatus = errors.New("invalid status") // retornado quando o status da fatura é inválido
)