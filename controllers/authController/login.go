package authcontroller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/Video-transcoder_Go/consts"
	"github.com/harsh082ip/Video-transcoder_Go/db"
	"github.com/harsh082ip/Video-transcoder_Go/helpers"
	"github.com/harsh082ip/Video-transcoder_Go/models"

	// "github.com/harsh082ip/URL-Shortener_Go/consts"
	// "github.com/harsh082ip/URL-Shortener_Go/db"
	// "github.com/harsh082ip/URL-Shortener_Go/helpers"
	// "github.com/harsh082ip/URL-Shortener_Go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Login(c *gin.Context) {

	var jsonData models.LoginUser
	// var apikey models.ApiKey
	// Bind and validate JSON
	var sessionInfo models.SessionInfo
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// Return a bad request response if there's an error in binding/validation
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error in request body",
			"error":  err.Error(),
		})
		return
	}

	hashPass := jsonData.Password
	collName := "Users"
	coll := db.OpenCollection(db.DBinstance(), collName)
	rdb := db.RedisConnect()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := coll.FindOne(ctx, bson.M{"email": jsonData.Email}).Decode(&jsonData)
	if err != nil {
		// Handle document not found error
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Error in Document",
				"error":  "No User found with the given details",
			})
			return
		}
		// Handle internal server error while searching for the user
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
			"error":  "Error is searching for user",
		})
		return
	}

	err = helpers.ComparePassword(jsonData.Password, hashPass)
	if err != nil {
		// Return unauthorized access response if passwords do not match
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized Access",
			"error":  "Password Mismatch",
		})
		return
	}

	sessionInfo.SessionID, err = helpers.CreateSeessionID(consts.SessionIDlength, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "SessionID generation error",
			"error":  err.Error(),
		})
		return
	}

	sessionInfo.Email = jsonData.Email

	// ----------------- SET to REDIS ----------------------------
	key := "session:" + sessionInfo.Email
	jsonSession, err := json.Marshal(sessionInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in Marshalling struct",
			"error":  err.Error(),
		})
		return
	}

	// first delete if older session exists
	rdb.Del(ctx, key)
	_, err = rdb.Set(ctx, key, jsonSession, consts.SessionTTL).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in Setting SessionID to DB",
			"error":  err.Error(),
		})
		return
	}

	collName = "SessionInfo"
	coll = db.OpenCollection(db.DBinstance(), collName)
	emailFilter := bson.M{"email": sessionInfo.Email}

	// Check if a document with the email already exists
	count, err := coll.CountDocuments(ctx, emailFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error in checking existing session",
			"error":  err.Error(),
		})
		return
	}
	// log.Println("here1...")

	if count == 0 {
		sessionInfo.CreatedAt = time.Now()
		sessionInfo.UpdatedAt = sessionInfo.CreatedAt
		// If no document with the email exists, insert a new document
		// log.Println("here2...")
		_, err := coll.InsertOne(ctx, sessionInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error in inserting new session",
				"error":  err.Error(),
			})
			rdb.Del(ctx, key)
			return
		}
	} else {
		// update the time
		// sessionInfo.UpdatedAt = time.Now()
		// log.Println("here3...")
		// If a document with the email exists, update the existing document
		update := bson.M{"$set": bson.M{"sessionid": sessionInfo.SessionID, "updatedat": time.Now()}}
		// log.Println("here3...")
		opts := options.Update().SetUpsert(false) // SetUpsert(false) for update only
		// log.Println("here3...")
		res, err := coll.UpdateOne(ctx, emailFilter, update, opts)
		// log.Println("here3...")
		log.Println(res.MatchedCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error in updating existing session",
				"error":  err.Error(),
			})
			rdb.Del(ctx, key)
			log.Println("here4....")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "User login Successful",
		"sessionInfo": sessionInfo,
		"count":       count,
		// "UpsertedCount": res.UpsertedCount,
		// "ModifiedCount": res.ModifiedCount,
	})
}
