package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"ui/html/pages/home.tpl.html"}, nil)

	checkErr(err)
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// simple validation
	
	if t.Title == "" {
		fmt.Println("excuting simple validation")	
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}

	// if input is okay, create a todo
	tm := todo{
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	fmt.Println("Created todomodel", tm)

	_, err := collection.InsertOne(context.Background(), tm)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
	return	
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		//"todo_id": tm.ID.Hex(),
	})
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	// Here's an array in which you can store the decoded documents
	var results []*todo

	// Finding multiple documents returns a cursor
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
		"message": "Failed to fetch todo",
		"error":   err,
		})
		return
	}

	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.Background()){

		// create a value into which the single document can be decoded
		var elem todo
		err := cursor.Decode(&elem)
		if err != nil {
			fmt.Println("Is it here")
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	defer cursor.Close(context.Background())
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": results,
	})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	fmt.Println(id)

	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}
	hid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": hid}

	deleteResult, err := collection.DeleteOne(context.Background(),filter)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
	return
	}
	fmt.Println(deleteResult)
	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo deleted successfully",})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}
	// simple validation
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}
	hid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": hid}
	update := bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}}
	_, err := collection.UpdateOne(context.Background(),filter,update)
	if err != nil {
		fmt.Println(err)
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})

}