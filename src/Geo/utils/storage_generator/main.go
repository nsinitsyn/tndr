package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Profile struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photos      string `json:"photos"`
}

func (p Profile) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func main() {
	var ctx = context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	profiles := map[string]interface{}{
		"11": Profile{11, "Alex", "Descr", `["123","456"]`},
		"12": Profile{12, "Kate", "Descr", `["123"]`},
		"13": Profile{13, "John", "Descr", `[]`},
	}

	// ucfe8: 55.481, 37.288
	err := client.HSet(ctx, "geohash:ucfe8", profiles).Err()
	logError(err)
}
