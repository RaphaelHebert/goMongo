package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/RaphaelHebert/goMongo/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userController struct{
	client *mongo.Client
}

func CreateNewUserController(c *mongo.Client) *userController {
	return &userController{c}
}

func (uc userController) UpdateUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	u := model.User{}

	// parse payload
	json.NewDecoder(req.Body).Decode(&u)
	// ad Id
	u.Id = primitive.NewObjectID()

	//insert to db
	res, err := uc.client.Database("tuto").Collection("users").UpdateByID(context.TODO(), u.Id, u)
	if err != nil {
		msg := fmt.Errorf("updateUser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return
	}
	fmt.Printf("Updated %v document with _id: %v\n", res.ModifiedCount, u.Id)

	// to send JSON
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "%v\n", u)
}

func (uc *userController) DeleteUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	u := model.User{}
	uid := p.ByName("id")
	// check if id is hex
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		msg := fmt.Errorf("getuser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return
	}
	u.Id = id
	// var res bson.M
	res, err := uc.client.Database("tuto").Collection("users").DeleteOne(context.TODO(), bson.D{{"_id", u.Id}})
	if err != nil {
		msg := fmt.Errorf("getuser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return
	}
	fmt.Printf("deleted %v user\n", res.DeletedCount)
	// to send JSON 
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(http.StatusOK)
	if res.DeletedCount < 1 {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%v\n",res.DeletedCount)
}


func (uc *userController) GetUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	u := model.User{}
	uid := p.ByName("id")
	// check if id is hex
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		msg := fmt.Errorf("getuser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return
	}

	u.Id = id
	var res bson.M
	err = uc.client.Database("tuto").Collection("users").FindOne(context.TODO(), bson.D{{"_id", u.Id}}).Decode(&res)
	if err != nil {
		msg := fmt.Errorf("getuser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return 
	}
	// to send JSON 
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n",res)
}

func (uc *userController) GetUsers(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	cur, err := uc.client.Database("tuto").Collection("users").Find(context.TODO(), bson.D{{}})
	if err != nil {
		msg := fmt.Errorf("getuser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return 
	}
	var results []model.User
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	resJSON, _ := json.Marshal(results)
	fmt.Fprintf(w, "%s\n", resJSON)
}

func (uc userController) CreateUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	u := model.User{}

	// get value from form
	u.Name = req.FormValue("name")
	u.Email = req.FormValue("email")

	// parse payload (overrides form value)
	json.NewDecoder(req.Body).Decode(&u)
	// ad Id
	u.Id = primitive.NewObjectID()

	//insert to db
	res, err := uc.client.Database("tuto").Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		msg := fmt.Errorf("createUser: %v", err)
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", msg)
		return
	}
	fmt.Printf("Inserted document with _id: %v\n", res.InsertedID)

	// to send JSON
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", u)
}
