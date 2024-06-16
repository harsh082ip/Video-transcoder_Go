package routes

import (
	"github.com/gin-gonic/gin"
	authcontroller "github.com/harsh082ip/Video-transcoder_Go/controllers/authController"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/auth/signup", authcontroller.SignUp)
	incomingRoutes.POST("/auth/login", authcontroller.Login)
}
