package s3controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	aws_conf "github.com/harsh082ip/Video-transcoder_Go/aws"
	"github.com/harsh082ip/Video-transcoder_Go/helpers"
	"github.com/harsh082ip/Video-transcoder_Go/models"
)

func PreSignedUrlToPutImage(c *gin.Context) {

	var fileinfo models.FileInfo
	if err := c.ShouldBindJSON(&fileinfo); err != nil {
		// Return a bad request response if there's an error in binding/validation
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error in request body, filename and email is required to proceed",
			"error":  err.Error(),
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	collName := "Users"
	res, _ := helpers.CheckIfDocExists("email", fileinfo.Email, collName, ctx)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "No user exists with the given email",
			"error":  "Please sent a correct email to proceed",
		})
		return
	}

	s3Client, err := aws_conf.GetS3Client()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "cannot initialize a s3 client",
			"error":  err.Error(),
		})
		return
	}
	// ctx, cancel :=

	uniqueKey, err := helpers.GetUniqueKey(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in creating unique key",
			"error":  err.Error(),
		})
		return
	}

	key := "videos/" + fileinfo.Email + "/" + uniqueKey + fileinfo.Filename
	input := &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(key),
		ContentType: aws.String("video/mp4"),
	}

	presignedClient := s3.NewPresignClient(s3Client)

	presignedURL, err := presignedClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in creating a pre-signed url :/",
			"error":  err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":        "Pre-Signed Url Creation Success",
		"URL":        presignedURL,
		"expires in": "15 Minutes",
	})
}
