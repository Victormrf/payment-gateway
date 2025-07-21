# --- Etapa de Build ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia os arquivos de módulos e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o restante do código
COPY . .

# Compila a aplicação com linkagem estática.
# CGO_ENABLED=0 é crucial para que o binário não dependa de bibliotecas de C
# e possa rodar na imagem scratch.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main cmd/app/main.go

# --- Etapa Final (Runtime) ---
FROM scratch

WORKDIR /app

# Copia o binário compilado da etapa de build
COPY --from=builder /app/main .

# Expõe a porta que sua aplicação usa
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["/app/main"]