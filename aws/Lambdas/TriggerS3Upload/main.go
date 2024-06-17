package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/harsh082ip/Video-transcoder_Go/db"
	"github.com/harsh082ip/Video-transcoder_Go/models"
)

func handleRequest(ctx context.Context, s3Event events.S3Event) {

	rdb := db.RedisConnect()
	var videojobs models.VideosJobs
	redis_key := "VideoJobs"
	ctx1, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println("Lambda Started")

	for _, record := range s3Event.Records {
		s3 := record.S3
		bucket := s3.Bucket.Name
		key := s3.Object.Key

		// Generate the S3 object URL
		// escapedKey := url.PathEscape(key)
		objectURL := fmt.Sprintf("https://s3.ap-south-1.amazonaws.com/%s/%s", bucket, key)

		videojobs.Key = key
		videojobs.ObjectUrl = objectURL

		jsonData, err := json.Marshal(videojobs)
		if err != nil {
			log.Fatal("Error in Marshalling the videojobs body", err.Error())
		}

		err = rdb.LPush(ctx1, redis_key, jsonData).Err()
		if err != nil {
			log.Fatal("Error in pushing VideoJobs to queue")
		}

		// Increase job count
		count, err := rdb.Get(ctx1, "Jobs").Result()
		if err != nil {
			log.Fatal("Error in Getting Jobs Count from redis", err.Error())
		}

		// Convert the value to an integer
		val, err := strconv.Atoi(count)
		if err != nil {
			log.Fatalf("Failed to Jobs value to integer: %v", err.Error())
		}

		// increse it by 1
		val++
		err = rdb.Set(ctx1, "Jobs", val, 0).Err()
		if err != nil {
			log.Fatal("Failed in Incresing Jobs Count, ", err.Error())
		}

		log.Println("Jobs pushed to queue, \nurl := ", videojobs.ObjectUrl)
	}
}

func main() {
	lambda.Start(handleRequest)
}
