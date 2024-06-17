package routes

import (
	"github.com/gin-gonic/gin"
	s3controller "github.com/harsh082ip/Video-transcoder_Go/controllers/s3Controller"
	"github.com/harsh082ip/Video-transcoder_Go/middleware"
)

func S3Routes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/videos/getpresignedurl", middleware.AuthMiddleware(), s3controller.PreSignedUrlToPutImage)
}
