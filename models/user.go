package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	Name string `json: "name", bson: "name"`
	Email string `json: "email", bson: "email"`
	Id primitive.ObjectID `json: "id", bson: "_id"`
}

type Users []User