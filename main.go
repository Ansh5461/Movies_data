package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"math/rand"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
	// * as we all know means a pointer
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

//we are passing a pointer of request sent by user in form of r, and w is response writer, so when we will send a response
//from this function outside, it will be w
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	//we will encode all the movies into w in form of json and send it to front end

}

//delete a movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, value := range movies {

		if value.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

//get a specific movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, value := range movies {
		if value.ID == params["id"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}
}

//create movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	//we will pass the arguments, and then encode them into w, which will be producing response
	w.Header().Set("Content-Type", "application/json")
	var movie Movie

	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

//update movie function- quite intimidating function, as being said
func updateMovie(w http.ResponseWriter, r *http.Request) {
	//approach, delete the movie which was already there, then add new movie
	w.Header().Set("Content_Type", "application/json")
	params := mux.Vars(r)

	for index, value := range movies {

		if value.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)

			movie.ID = params["id"]

			movies = append(movies, movie)

			json.NewEncoder(w).Encode(movie)
			return
		}

	}

}

func main() {

	r := mux.NewRouter()

	//lets create some movies which will already be present to just check whether our router is working on all fronts or not

	//we want the referennce of the address of director &-give the address, *- will access the address
	movies = append(movies, Movie{ID: "1", Isbn: "438794", Title: "Spiderman - No way home", Director: &Director{Firstname: "Jon", Lastname: "Watts"}})
	movies = append(movies, Movie{ID: "2", Isbn: "696941", Title: "Moon knight", Director: &Director{Firstname: "Mohammed", Lastname: "Diab"}})

	//here we will create 5 crud operations, using handle func

	//we will have 5 different routes and 5 different functions, which are getall, getbyid, create, update, delete - these are 5 routes
	//and we have functions associated with them too

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
