package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jobscontroller "github.com/harsh082ip/Video-transcoder_Go/controllers/jobsController"
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
		err = jobscontroller.IncreaseJobsCount(ctx1)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Jobs pushed to queue, \nurl := ", videojobs.ObjectUrl)
	}
}

func main() {
	lambda.Start(handleRequest)
}
