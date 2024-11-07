package models

import "time"

type Form struct {
	Email string    `json:"email" bson:"email" binding:"email,required"`
	Id    string    `json:"id" bson:"id" binding:"required"`
	Ts    time.Time `json:"ts" bson:"ts" binding:"required"`
	Data  any       `json:"data" bson:"data" binding:"required"`
}
