package main
import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	
)

type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string            `json: "title,omitempty" bson:"title,omitempty"`
	Genre  string            `json: "genre,omitempty" bson:"genre,omitempty"`
	Author  *Author          `json:"author" bson:"author,omitempty"`	
}

type Author struct {
	FirstName string        `json: "FirstName,omitempty" bson: "FirstName,omitempty"`
	LastName string         `json: "LastName,omitempty" bson: "LastName,omitempty"`
}

func ConnectToDB() *mongo.Collection{

   clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

   client, err := mongo.Connect(context.TODO(), clientOptions)

   if err != nil{
	   log.Fatal("An Error occured while connecting the DB")
   }

   fmt.Println("MongoDB Connected")

   collection := client.Database("test_book_db").Collection("books")

   return collection
}


var collection = ConnectToDB()

 func getBooks (w http.ResponseWriter, r *http.Request ){
   w.Header().Set("Content-Type", "application/json")
  var books []Book
  cur, err := collection.Find(context.TODO(), bson.M{})
 
  defer cur.Close(context.TODO())

  for cur.Next(context.TODO()){

	  var book Book

	err :=  cur.Decode(&book)
	
	if err != nil{
		log.Fatal(err)
	}

   books = append(books, book)
  }
  if err != nil{
    log.Fatal(err)
  }

  json.NewEncoder(w).Encode(books)
 }
 func getBook(w http.ResponseWriter, r *http.Request){


  w.Header().Set("Content-type", "application/json")
  var book Book

  params := mux.Vars(r)

  _id, _ := primitive.ObjectIDFromHex(params["id"])

  filter := bson.M{"_id":_id}

  err := collection.FindOne(context.TODO(), filter).Decode(&book)
  if err != nil{
    log.Fatal(err)
  }
  json.NewEncoder(w).Encode(book)


 }

 func createBook(w http.ResponseWriter, r *http.Request){

	 w.Header().Set("Content-Type", "application/json")
	 var book Book

	 json.NewDecoder(r.Body).Decode(&book)

	 result, err := collection.InsertOne(context.TODO(), book)

	 if err != nil{
		log.Fatal(err)
	 }
	 json.NewEncoder(w).Encode(result)
 }
 func editBook(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	var book Book
	params := mux.Vars(r)
	_id, _ := primitive.ObjectIDFromHex(params["id"])

	json.NewDecoder(r.Body).Decode(&book)

	filter := bson.M{"_id": _id}

	update := bson.D{
		{
			"$set", bson.D{
				{"title", book.Title},
				{"genre", book.Genre},
				{"author", bson.D{
				  {"FirstName", book.Author.FirstName},
				  {"LastName", book.Author.LastName},
				}},
			},
		},
	}

	
	 err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)
	
	

	if err != nil{
	   log.Fatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request){


	w.Header().Set("Content-type", "application/json")
  
	params := mux.Vars(r)
  
	_id, _ := primitive.ObjectIDFromHex(params["id"])
  
	filter := bson.M{"_id":_id}
  
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil{
	  log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
  
  
   }
func main(){
	r := mux.NewRouter()
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", editBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", r))

}