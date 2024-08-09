.DEFAULT_GOAL = build

build:
	@go mod tidy
	@go build -o bin/echoes ./main.go
	@cp config.yaml bin

run: build
	@ECHOES_ADDR=0.0.0.0:5000 ./bin/echoes
