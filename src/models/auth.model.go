package models

type Token struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Token    string `json:"token" bson:"token" binding:"required"`
}
