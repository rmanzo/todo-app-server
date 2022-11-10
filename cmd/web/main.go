package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var rnd *renderer.Render

type todo struct {
		ID        	string    	`json:"id" bson:"_id,omitempty"`
		Title     	string    	`json:"title" bson:"title"`
		Completed 	bool      	`json:"completed" bson:"completed"`
		CreatedAt 	time.Time 	`json:"createdAt" bson:"createAt"`
		CompletedAt time.Time   `json:"completedAt" bson:"completedAt"`
	}

func main() {

    hostName := flag.String("hostName", "mongodb://localhost:27017", "Mongo DB host")
	dbName := flag.String("dbName", "demo_todo", "Mongo DB Name")
	collectionName := flag.String("collectionName", "todo", "Mongo Collection Name")
	port := flag.String("port", ":8000", "HTTP network address")

    flag.Parse()

	rnd = renderer.New()
	//client option
	clientOptions := options.Client().ApplyURI(*hostName)

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("unable to connect to MongoDB, please check it is running")
	}
	fmt.Println("Connected to MongoDB")

	//collection reference
	collection = client.Database(*dbName).Collection(*collectionName)
	fmt.Println("Collection reference is ready")


	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         *port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening on port ", *port)
	if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
	}	
	
}