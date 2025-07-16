package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)


type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu  	sync.RWMutex // Bloqueia a escrita concorrente de valor
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateAPIKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewAccount(name, email string) *Account {
	account := &Account{
		ID: uuid.New().String(),
		Name:      name,
		Email:     email,
		Balance:  0.0,
		APIKey:  generateAPIKey(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

// criação de uma função que vai atuar como um "método" da classe Account
func (a *Account) AddBalance(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock() // O defer vai rodar sempre por ultimo
	a.Balance += amount
	a.UpdatedAt = time.Now()
}