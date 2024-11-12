package models

type User struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"pw" bson:"pw" binding:"required"`
}
