package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../errors"
	"../models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetImages(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// we created Image array
	var images []models.Image
	collection := db.Collection("images")
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		errors.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var image models.Image
		// & character returns the memory address of the following variable.
		err := cur.Decode(&image) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		images = append(images, image)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(images) // encode similar to serialize process.
}

func GetImage(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var image models.Image
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := db.Collection("images")

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&image)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(image)
}

func CreateImage(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var image models.Image
	_ = json.NewDecoder(r.Body).Decode(&image)
	collection := db.Collection("images")
	fmt.Println(image)
	result, err := collection.InsertOne(context.TODO(), image)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func UpdateImage(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var image models.Image

	collection := db.Collection("images")

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&image)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"isbn", image.Isbn},
			{"title", image.Title},
			{"author", bson.D{
				{"firstname", image.Author.FirstName},
				{"lastname", image.Author.LastName},
			}},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&image)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	image.ID = id

	json.NewEncoder(w).Encode(image)
}

func DeleteImage(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}
	collection := db.Collection("images")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
