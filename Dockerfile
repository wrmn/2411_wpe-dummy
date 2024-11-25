FROM golang:alpine

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest && go mod download

CMD ["air", "-c", ".air.toml"]