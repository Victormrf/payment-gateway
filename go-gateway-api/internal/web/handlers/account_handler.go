package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/dto"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/service"
)

type AccountHandler struct {
	accountService *service.AccountService // Este handler vai precisar acessar o service, logo o service vai ser uma dependencia dele
}

func NewAccountHandlers(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// Endpoints que o handler vai expor (controller com requests e responses)
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAccountInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.accountService.CreateAccount(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API Key is required", http.StatusUnauthorized)
		return
	}

	output, err := h.accountService.FindByAPIKey(apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}