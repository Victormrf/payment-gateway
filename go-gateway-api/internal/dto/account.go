package dto

import (
	"time"

	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/domain"
)

type CreateAccountInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	APIKey    string    `json:"api_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Quando eu tenho um DTO e quero transformar ele em um objeto de domínio
func ToAccount(input CreateAccountInput) *domain.Account {
	return domain.NewAccount(input.Name, input.Email)
}

// Quando eu tenho um objeto de domínio e quero transformar ele em um DTO
func FromAccount(account *domain.Account) AccountOutput {
	return AccountOutput{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Balance:   account.Balance,
		APIKey:    account.APIKey,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}