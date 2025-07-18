package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver de migração para PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Fonte de migração de arquivos

	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/repository"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/service"
	"github.com/Victormrf/payment-gateway/go-gateway-api/internal/web/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}


// Onde vou juntar Web e Application
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// string de conexão com o banco de dados
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	// URL de conexão para a ferramenta de migração
	migrateURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		log.Printf("Tentando aplicar migrações... Tentativa %d/%d", i+1, maxRetries)
		m, err := migrate.New(
			"file://migrations", // Certifique-se de que este é o caminho correto para seus arquivos de migração
			migrateURL,
		)
		if err != nil {
			log.Printf("Erro ao criar instância de migração: %v. Retentando em 5 segundos...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			// Se o erro for de conexão, tentamos novamente
			if err.Error() == "dial tcp: lookup db: no such host" || err.Error() == "pq: SSL is not enabled on the server" { // Exemplo de erros de conexão
				log.Printf("Erro de conexão com o banco de dados durante a migração: %v. Retentando em 5 segundos...", err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Fatalf("Falha ao executar as migrações: %v", err)
		} else if err == migrate.ErrNoChange {
			log.Println("Nenhuma nova migração para aplicar.")
		} else {
			log.Println("Migrações aplicadas com sucesso!")
		}
		break // Saia do loop se as migrações foram bem-sucedidas ou não houveram mudanças
	}


	// Cofiguração de conexão com banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, *accountService)

	port := getEnv("HTTP_PORT", "8080")
	srv := server.NewServer(accountService, invoiceService, port)
	srv.ConfigureRoutes()

	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
