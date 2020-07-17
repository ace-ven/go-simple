package services

import (
	"context"
	"encoding/json"
	"net/http"

	"../errors"
	"../models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InserUser(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	userCollection := db.Collection("users")

	results, err := userCollection.InsertOne(context.TODO(), user)

	if err != nil {
		errors.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(results)

}
