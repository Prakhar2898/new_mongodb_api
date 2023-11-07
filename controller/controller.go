package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Prakhar2898/mongoapi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://dbPrakhar:mongodbgoprakhar@cluster0.kocmpra.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const collName = "watchlist"

// imp
var col *mongo.Collection

// connect with mongoDB
func init() {
	//client option
	connectOption := options.Client().ApplyURI(connectionString)

	//connect with mongodb
	client, err := mongo.Connect(context.TODO(), connectOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo DB connect succ...")

	col = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance is ready...")
}

// MONGODB helpers - file

// insert 1 record
func insertOneMovie(movie model.Netflix) {
	inserted, err := col.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted One Movie to DB with Id: ", inserted)

}

// update one
func updateOneMovie(mId string) {
	id, _ := primitive.ObjectIDFromHex(mId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified count: ", result.ModifiedCount)

}
func errHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}

	deleteResult, err := col.DeleteOne(context.Background(), filter)
	errHandle(err)

	fmt.Println("delete count: ", deleteResult.DeletedCount)

}

func deleteAllMovies() int64 {
	deleteResult, err := col.DeleteMany(context.Background(), bson.D{{}}, nil)
	errHandle(err)

	fmt.Println("Delete count(delete all): ", deleteResult)
	return deleteResult.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := col.Find(context.Background(), bson.D{{}})
	errHandle(err)

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		errHandle(err)
		movies = append(movies, movie)
	}
	return movies
}

// Actual controller - file

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlecode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(&movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	para := mux.Vars(r)
	updateOneMovie(para["id"])
	json.NewEncoder(w).Encode(para["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	v := deleteAllMovies()
	json.NewEncoder(w).Encode(v)

}
