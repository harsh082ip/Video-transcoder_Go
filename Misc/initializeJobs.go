package misc

import (
	"context"
	"log"
	"time"

	"github.com/harsh082ip/Video-transcoder_Go/db"
)

func InitializeJobs() {
	rdb := db.RedisConnect()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	key := "Jobs"

	// Set the key to "0" if it does not already exist
	set, err := rdb.SetNX(ctx, key, "0", 0).Result()
	if err != nil {
		log.Fatal("Initializing jobs error", err.Error())
	}

	if set {
		log.Println("Key did not exist, JOBS initialized to 0")
	} else {
		log.Println("JOSB already exists, no action taken")
	}
}
