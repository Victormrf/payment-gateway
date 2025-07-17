package server

import (
	"net/http"

	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/service"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/web/handlers"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux // roteador do go
	server *http.Server
	accountService *service.AccountService
	port string
}

func NewServer(accountService *service.AccountService, port string) *Server {
	return &Server{
		router: chi.NewRouter(),
		accountService: accountService,
		port: port,
	}
}

func (s *Server) ConfigureRoutes() {
	accountHandler := handlers.NewAccountHandlers(s.accountService)

	s.router.Post("/accounts", accountHandler.Create)
	s.router.Get("/accounts", accountHandler.Get)
} 

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr: ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}