package app

import (
	"log"
	"net/http"

	"./dbconnector"
	"./services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize() {
	a.DB = dbconnector.ConnectMongoDB()
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Post("/api/users", a.handleRequest(services.InserUser))
	a.Get("/api/users/findByName/{name}", a.handleRequest(services.FindUserByName))
	a.Get("/api/images", a.handleRequest(services.GetImages))
	a.Get("/api/images/{id}", a.handleRequest(services.GetImage))
	a.Post("/api/images", a.handleRequest(services.CreateImage))
	a.Put("/api/images/{id}", a.handleRequest(services.UpdateImage))
	a.Delete("/api/images/{id}", a.handleRequest(services.DeleteImage))
	a.Post("/auth/login", a.handleRequest(services.Signin))

}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

type RequestHandlerFunction func(db *mongo.Database, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(services RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services(a.DB, w, r)
	}
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
