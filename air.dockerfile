FROM golang:1.22-alpine3.20

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

RUN apk --update --no-cache add make

CMD ["air", "-c", ".air.toml"]