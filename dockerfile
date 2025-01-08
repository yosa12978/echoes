FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN apk add --no-cache coreutils
RUN make 

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/config.yaml .

RUN apk --update --no-cache add curl

EXPOSE 80
ENTRYPOINT ["./echoes"]
