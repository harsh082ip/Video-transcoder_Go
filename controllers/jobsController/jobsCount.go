package jobscontroller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/harsh082ip/Video-transcoder_Go/db"
)

func IncreaseJobsCount(ctx context.Context) error {

	rdb := db.RedisConnect()

	// Increase job count
	count, err := rdb.Get(ctx, "Jobs").Result()
	if err != nil {
		return fmt.Errorf("error in Getting Jobs Count from redis, %v", err.Error())
	}

	// Convert the value to an integer
	val, err := strconv.Atoi(count)
	if err != nil {
		return fmt.Errorf("failed to Jobs value to integer: %v", err.Error())
	}

	// increse it by 1
	val++
	err = rdb.Set(ctx, "Jobs", val, 0).Err()
	if err != nil {
		return fmt.Errorf("failed in Incresing Jobs Count, %v", err.Error())
	}

	return nil
}

func DecreaseJobsCount(ctx context.Context) error {

	rdb := db.RedisConnect()

	// Decrease job count
	count, err := rdb.Get(ctx, "Jobs").Result()
	if err != nil {
		return fmt.Errorf("error in Getting Jobs Count from redis, %v", err.Error())
	}

	// Convert the value to an integer
	val, err := strconv.Atoi(count)
	if err != nil {
		return fmt.Errorf("failed to Jobs value to integer: %v", err.Error())
	}

	// decrease it by 1
	val--
	err = rdb.Set(ctx, "Jobs", val, 0).Err()
	if err != nil {
		return fmt.Errorf("failed in Incresing Jobs Count, %v", err.Error())
	}

	return nil
}
