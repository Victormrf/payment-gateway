package service

import (
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/domain"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/dto"
)

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

func (s *AccountService) CreateAccount(input dto.CreateAccountInput) (*dto.AccountOutput, error) {
	account := dto.ToAccount(input)
	existingAccount, err := s.repository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err 
	}

	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey // Se já existir uma conta com a mesma API Key, retorna erro
	}

	err = s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil // Retorna o DTO da conta criada
}

func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil // Retorna o DTO da conta criada
}

func (s *AccountService) FindByAPIKey(apiKey string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil // Retorna o DTO da conta encontrada
}

func (s *AccountService) FindByID(id string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil // Retorna o DTO da conta encontrada
}