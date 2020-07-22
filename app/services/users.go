package services

import (
	"context"
	"encoding/json"
	"net/http"

	"../errors"
	"../models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InserUser(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	password := []byte(user.Password)
	passHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	user.Password = string(passHash)

	userCollection := db.Collection("users")
	results, err := userCollection.InsertOne(context.TODO(), user)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(results)

}

func GetUser(db *mongo.Database, name string) models.User {
	var user models.User
	collection := db.Collection("users")
	filter := bson.M{"firstname": name}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		panic(err)
	}
	return user

}

func FindUserByName(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	firstname, _ := params["name"]
	user := GetUser(db, firstname)

	json.NewEncoder(w).Encode(user)
}
