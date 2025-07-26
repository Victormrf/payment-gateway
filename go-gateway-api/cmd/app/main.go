package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

// Onde vou juntar web e application
func main() {
	if getEnv("APP_ENV", "development") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file for development: ", err)
		}
		log.Println("Loaded .env file for development environment.")
	} else {
		log.Println("Running in production environment, using system environment variables.")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

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
			"file://migrations",
			migrateURL,
		)
		if err != nil {
			log.Printf("Erro ao criar instância de migração: %v. Retentando em 5 segundos...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			if isNetworkError(err) {
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
		break 
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database after migrations: %v", err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// Configura e inicializa o Kafka
	baseKafkaConfig := service.NewKafkaConfig()

	// Configura e inicializa o produtor Kafka
	producerTopic := getEnv("KAFKA_PRODUCER_TOPIC", "pending_transactions")
	producerConfig := baseKafkaConfig.WithTopic(producerTopic)
	kafkaProducer := service.NewKafkaProducer(producerConfig)
	defer kafkaProducer.Close()

	// Inicializa camadas da aplicação (repository -> service -> server)
	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, *accountService, kafkaProducer)

	// Configura e inicializa o consumidor Kafka
	consumerTopic := getEnv("KAFKA_CONSUMER_TOPIC", "transaction_results")
	consumerConfig := baseKafkaConfig.WithTopic(consumerTopic)
	groupID := getEnv("KAFKA_CONSUMER_GROUP_ID", "gateway-group")
	kafkaConsumer := service.NewKafkaConsumer(consumerConfig, groupID, invoiceService)
	defer kafkaConsumer.Close()

	// Inicia o consumidor Kafka em uma goroutine
	go func() {
		if err := kafkaConsumer.Consume(context.Background()); err != nil {
			log.Printf("Error consuming kafka messages: %v", err)
		}
	}()

	port := getEnv("HTTP_PORT", "8080")
	srv := server.NewServer(accountService, invoiceService, port)
	srv.ConfigureRoutes()

	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

// isNetworkError tenta identificar erros de rede comuns
func isNetworkError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "no such host") || // Docker pode não resolver 'db' imediatamente
		strings.Contains(errMsg, "connection refused") || // DB não está pronto ou inacessível
		strings.Contains(errMsg, "timeout") // Timeout de conexão
}
