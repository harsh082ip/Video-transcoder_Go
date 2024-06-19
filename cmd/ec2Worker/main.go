package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	ecscontroller "github.com/harsh082ip/Video-transcoder_Go/controllers/ecsController"
	jobscontroller "github.com/harsh082ip/Video-transcoder_Go/controllers/jobsController"
	"github.com/harsh082ip/Video-transcoder_Go/db"
	ecshelper "github.com/harsh082ip/Video-transcoder_Go/helpers/ecsHelper"
	"github.com/harsh082ip/Video-transcoder_Go/models"
	"github.com/redis/go-redis/v9"
)

func main() {

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	var videojobs models.VideosJobs
	rdb := db.RedisConnect()
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	DESTINATION_BUCKET_NAME := os.Getenv("DESTINATION_BUCKET_NAME")

	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" || DESTINATION_BUCKET_NAME == "" {
		log.Fatal("Some required env are needed to proceed")
	}

	// START PROCESS
	for {

		res, err := rdb.RPop(ctx, "videos").Result()
		if err == redis.Nil {
			// No jobs found, check again
			log.Println("Checking Jobs again")
			time.Sleep(time.Second * 2)
			continue
		}
		if err != nil {
			log.Fatal("error in poping the jobs from the queue, ", err.Error())
		}
		err = json.Unmarshal([]byte(res), &videojobs)
		if err != nil {
			log.Fatal("error in unmarshalling jobs body", err.Error())
		}
		for {
			runningTasks, err := ecshelper.ListRunningTask()
			if err != nil {
				log.Fatal("error in listing running tasks, ", err.Error())
			}
			if runningTasks >= 5 {
				time.Sleep(time.Second * 2)
				continue
			}
			DESTINATION_1080 := "s3://" + DESTINATION_BUCKET_NAME + "/" + videojobs.Key + "1080.mp4"
			DESTINATION_720 := "s3://" + DESTINATION_BUCKET_NAME + "/" + videojobs.Key + "720.mp4"
			DESTINATION_360 := "s3://" + DESTINATION_BUCKET_NAME + "/" + videojobs.Key + "360.mp4"

			resp, err := ecscontroller.RunECSTask(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, videojobs.ObjectUrl, DESTINATION_1080, DESTINATION_720, DESTINATION_360)
			if err != nil || !resp {
				log.Fatal("error in running tasks, ", err.Error())
			}

			err = jobscontroller.DecreaseJobsCount(ctx)
			if err != nil {
				log.Fatal("Error in Decreasing Job Count")
			}
			break
		}
	}
}
