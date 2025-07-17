package server

import (
	"net/http"

	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/service"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/web/handlers"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux // roteador do go
	server *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port string
}

func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService, port string) *Server {
	return &Server{
		router: chi.NewRouter(),
		accountService: accountService,
		invoiceService: invoiceService,
		port: port,
	}
}

func (s *Server) ConfigureRoutes() {
	accountHandler := handlers.NewAccountHandlers(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)

	s.router.Post("/accounts", accountHandler.Create)
	s.router.Get("/accounts", accountHandler.Get)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		s.router.Post("/invoice", invoiceHandler.Create)
		s.router.Get("/invoice/{id}", invoiceHandler.GetByID)
		s.router.Get("/invoice", invoiceHandler.ListByAccount)
	})

} 

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr: ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}