package main

import (
	"os"
	"os/signal"
	"syscall"
	"tinder-geo/internal/app"
)

// $env:CONFIG_PATH = '././config/config.yaml'; $env:CGO_ENABLED=1; $env:CC="C:\TDM-GCC-64\bin\gcc"; go run cmd/tinder-geo/main.go

func main() {

	// fmt.Println("!")

	// g := &run.Group{}
	// g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	// ctx, cancel := context.WithCancel(context.Background())
	// g.Add(func() error {
	// 	defer cancel()
	// 	select {
	// 	case <-time.After(15 * time.Second):
	// 	case <-ctx.Done():
	// 		break
	// 	}
	// 	return nil
	// },
	// 	func(err error) { cancel() })

	// if err := g.Run(); err != nil {
	// 	log.Fatal(err)
	// }

	closer := app.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	closer()
}
