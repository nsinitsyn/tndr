package main

import (
	"os"
	"os/signal"
	"syscall"
	"tinder-geo/internal/app"
)

// $env:CONFIG_PATH = '././config/config.yaml'; go run cmd/tinder-geo/main.go

// with kafka:
// $env:CONFIG_PATH = '././config/config.yaml'; $env:CGO_ENABLED=1; $env:CC="C:\TDM-GCC-64\bin\gcc"; go run cmd/tinder-geo/main.go

// kafka client readme:
// https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#getting-started

// tdm-gcc:
// https://github.com/mattn/go-sqlite3/issues/168#issuecomment-1528722456
// https://github.com/jmeubank/tdm-gcc
// https://jmeubank.github.io/tdm-gcc/download/

func main() {
	closer := app.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	closer()

	// match service http client - 100 ошибок go - как закрывать body http клиентов правильно
	// observability - 2 видео art of development
	// tests
	// retry и circuit breaker к внешним системам

	// Пересмотреть где для сервисов использовать указатели, а где нет. Лучше по возможности все передавать и отдавать по значению, особенно возвращать интерфейсы из функций - переделать.

	// Валидация grpc запросов
	// Настроить таймауты в grpc сервере и keepAlive

	// feature: remove profile
}
