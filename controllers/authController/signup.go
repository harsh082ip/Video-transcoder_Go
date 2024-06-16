package authcontroller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/Video-transcoder_Go/db"
	"github.com/harsh082ip/Video-transcoder_Go/helpers"
	"github.com/harsh082ip/Video-transcoder_Go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func SignUp(c *gin.Context) {

	var jsonData models.User

	// Bind and validate JSON
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// Return a bad request response if there's an error in binding/validation
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error in request body",
			"error":  err.Error(),
		})
		return
	}

	collName := "Users"
	coll := db.OpenCollection(db.DBinstance(), collName)

	emailFilter := bson.M{"email": jsonData.Email}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	count, err := coll.CountDocuments(ctx, emailFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error While Checking for Doc",
			"error":  err.Error(),
			"count":  count,
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Doc Duplication not allowed",
			"error":  "this email already exists",
		})
		return
	}

	jsonData.Password, err = helpers.HashPassword(jsonData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in generating hash for password",
			"error":  err.Error(),
		})
		return
	}

	jsonData.ID = primitive.NewObjectID()

	// finally add the user to db
	_, err = coll.InsertOne(ctx, jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error is attempting to SignUp",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "User SignUp Successful",
	})
}
