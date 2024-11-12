package main

import (
	"log"
	"tinder-geo/internal/app"
)

// $env:CONFIG_PATH = '././config/config.yaml'; go run cmd/tinder-geo/main.go

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// code := geohash.EncodeWithPrecision(55.7893, 37.7717, 5)
	// fmt.Println(code)
	// // clean arh: app в коде Олега Козырева и видео art of development
	// // grpc server
	// // kafka client
	// // mongo client
	// // match service http client - 100 ошибок go - как закрывать body http клиентов правильно
	// // business logic - searching nearby profiles by geohash
	// // observability - 2 видео art of development
	// // graceful shutdown
	// // tests
}
