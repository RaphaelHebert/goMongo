package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/RaphaelHebert/goMongo/controllers"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tpl *template.Template
var client *mongo.Client
var err error

// TODO: make an utils function errorInfo
// TODO: pass error to login if db cannot init
// TODO: mv getClient to utils
// TODO: add comments for go doc 
// TODO: add production host and port

func init(){
	// init templates
	tpl = template.Must(template.ParseGlob("templates/*"))

	// init db
	// allow fails, client must be able to see UI and error message
	client = getClient()
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		msg := fmt.Errorf("could not ping db: %s", err)
		fmt.Println(msg)
	}
}

func main() {
	pc := controllers.CreateNewPageController(tpl)
	uctl := controllers.CreateNewUserController(client)

	jr := httprouter.New()
	jr.GET("/", pc.Index)
	jr.DELETE("/user/:id", uctl.DeleteUser)
	jr.GET("/user/:id", uctl.GetUser)
	jr.POST("/user", uctl.CreateUser)
	jr.PUT("/user", uctl.UpdateUser)
	jr.GET("/users", uctl.GetUsers)
	http.ListenAndServe("localhost:8080", jr)
	
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func getClient() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}
