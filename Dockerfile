# STAGE 1: Build a aplicação
# Usamos um "alias" 'builder' para esta fase.
FROM golang:1.21-alpine AS builder

# É uma boa prática rodar as ferramentas como um usuário não-root
RUN adduser -D -g '' appuser

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências primeiro
# O Docker armazena em cache esta camada. Ele só irá refazer o download
# das dependências se o go.mod ou go.sum mudar.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte do projeto
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 cria um binário estaticamente linkado, essencial para imagens mínimas.
# -o /url-shortener define o nome e local do arquivo de saída.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /url-shortener ./cmd/api/main.go

# STAGE 2: Cria a imagem final e leve
# Começamos de uma imagem base mínima. 'alpine' é pequena e segura.
FROM alpine:latest

# Adiciona o mesmo usuário não-root da fase de build
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Define o diretório de trabalho
WORKDIR /app

# Copia apenas o binário compilado da fase 'builder'.
# Isso mantém a imagem final extremamente pequena e segura, sem código-fonte ou ferramentas de build.
COPY --from=builder /url-shortener .

# A aplicação cria o diretório 'data' sozinha, mas aqui garantimos que o
# diretório /app/data pertencerá ao nosso usuário, para que a aplicação tenha permissão para escrever nele.
# É nesta pasta que montaremos nosso Volume para persistir os dados.
RUN mkdir data && chown -R appuser:appgroup /app/data

# Troca para o usuário não-root para rodar a aplicação
USER appuser

# Expõe a porta 8080 para o mundo exterior (fora do contêiner)
EXPOSE 8080

# Comando para executar a aplicação quando o contêiner iniciar
CMD ["./url-shortener"]
