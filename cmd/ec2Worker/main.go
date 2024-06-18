package main

import (
	"log"

	ecshelper "github.com/harsh082ip/Video-transcoder_Go/helpers/ecsHelper"
)

func main() {

	for {
		count, err := ecshelper.ListRunningTask()
		if err != nil {
			log.Fatal(err.Error())
		}
		if count == 0 {
			log.Println("Checking again")
			continue
		}
	}
}
