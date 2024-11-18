package main

import (
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
	app.Run()

	// code := geohash.EncodeWithPrecision(55.7893, 37.7717, 5)
	// fmt.Println(code)
	// // clean arh: app в коде Олега Козырева и видео art of development
	// // grpc server
	// // kafka client
	// // mongo client
	// // match service http client - 100 ошибок go - как закрывать body http клиентов правильно
	// // business logic - searching nearby profiles by geohash
	// // observability - 2 видео art of development
	// // tests
	// // alive, ready endpoints
	// // retry и circuit breaker к внешним системам

	// todo: messaging должен обращаться к сервису, а не к стораджу!

	// Пересмотреть где для сервисов использовать указатели, а где нет. Лучше по возможности все передавать и отдавать по значению, особенно возвращать интерфейсы из функций - переделать.

	// Логирование в файл
	// Валидация grpc запросов
	// Настроить таймауты в grpc сервере

	// feature: remove profile
}
