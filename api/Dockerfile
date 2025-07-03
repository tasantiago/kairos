FROM golang:1.21

WORKDIR /app

# Instala as dependências
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Mantém o container "aberto" para desenvolvimento
CMD ["tail", "-f", "/dev/null"]
