package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	misc "github.com/harsh082ip/Video-transcoder_Go/Misc"
	"github.com/harsh082ip/Video-transcoder_Go/consts"
	"github.com/harsh082ip/Video-transcoder_Go/routes"
)

func main() {

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.S3Routes(router)

	// Initialize Jobs
	go misc.InitializeJobs()

	if err := http.ListenAndServe(consts.WEBPORT, router); err != nil {
		log.Fatal("Error starting the server on", consts.WEBPORT, err.Error())
	}
}
