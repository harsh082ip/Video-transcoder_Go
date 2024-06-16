package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/Video-transcoder_Go/consts"
	"github.com/harsh082ip/Video-transcoder_Go/routes"
)

func main() {

	router := gin.Default()

	routes.AuthRoutes(router)

	if err := http.ListenAndServe(consts.WEBPORT, router); err != nil {
		log.Fatal("Error starting the server on", consts.WEBPORT, err.Error())
	}
}
