package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `json:"name" binding:"required"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password" binding:"required,min=6"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SessionInfo struct {
	Email     string    `json:"email" binding:"required,email"`
	SessionID string    `json:"sessionid" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedNow"`
}

type FileInfo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `json:"email" binding:"required"`
	Filename  string             `json:"filename" binding:"required"`
	S3Address string             `json:"s3address"`
}

type VideosJobs struct {
	Key       string `json:"key"`
	ObjectUrl string `json:"object_url"`
}
