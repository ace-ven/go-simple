package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SigninUser struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func Signin(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	password := []byte(user.Password)

	if err != nil {
		println(err)
	}
	dbUser := GetUser(db, user.FirstName)
	var hash = []byte(dbUser.Password)
	err = bcrypt.CompareHashAndPassword(hash, password)

	if err != nil {
		fmt.Println(err)
	}
	token, err := CreateToken(dbUser.ID)

	newSignin := SigninUser{Name: dbUser.FirstName, Token: token}
	u, err := json.Marshal(newSignin)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(u))

}

func CreateToken(userid primitive.ObjectID) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
