package main

// Importing the necessary packages
import (
	"fmt"          // For formatted I/O operations
	"log"          // For logging errors
	"encoding/json" // For JSON encoding and decoding
	"math/rand"    // 
	"net/http"     // For HTTP server functionalities
	"strconv"      // 
	"github.com/gorilla/mux" // For HTTP request multiplexing
)

// Movie struct to define the shape of Movie objects
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director struct to define the shape of Director objects
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Variable to hold a slice of Movie structs
var movies []Movie

// Handler to get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode the movies slice to JSON and send it
	json.NewEncoder(w).Encode(movies)
}

// Handler to delete a movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Get URL parameters
	params := mux.Vars(r)

	// Loop through movies to find the one to delete
	for index, movie := range movies {
		if movie.ID == params["id"] {
			// Delete the movie from the slice
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	// Encode the updated movies slice to JSON and send it
	json.NewEncoder(w).Encode(movies)
}

// Handler to get a single movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Get URL parameters
	params := mux.Vars(r)

	// Loop through movies to find the one with the specified ID
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	// If movie not found, return empty Movie struct
	json.NewEncoder(w).Encode(&Movie{})
}

// createMovie is a HTTP handler function that creates a new movie.
func createMovie(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON.
	// This informs the client that the server will be returning JSON-formatted data.
	w.Header().Set("Content-Type", "application/json")
	
	// Declare a new Movie struct variable.
	// This will hold the incoming JSON payload from the client request.
	var movie Movie
	
	// Decode the JSON body of the incoming HTTP request into the movie variable.
	// Notice the use of `&movie`, passing a pointer to the movie variable.
	// This allows json.NewDecoder to populate the movie variable directly.
	_ = json.NewDecoder(r.Body).Decode(&movie)  // Error ignored for brevity, but it's good to handle it.
	
	// Generate a random ID for the new movie.
	// The rand.Intn function returns a random integer n, where 0 <= n < 1000000.
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	
	// Append the new movie to the global movies slice.
	// The append function returns a new slice containing all the items.
	movies = append(movies, movie)
	
	// Encode the newly created movie to JSON and write it to the response.
	// This gives the client confirmation of what was created.
	json.NewEncoder(w).Encode(movie)
}

// updateMovie is a HTTP handler function that updates an existing movie.
func updateMovie(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON.
	// This informs the client that the server will be returning JSON-formatted data.
	w.Header().Set("Content-Type", "application/json")

	// Fetch URL parameters using Gorilla Mux's Vars function.
	// It returns a map of the URL variables where the key is the variable name and the value is its value.
	params := mux.Vars(r)

	// Loop through the global movies slice to find the movie to be updated.
	for index, movie := range movies {
		// If the movie's ID matches the ID from the URL parameter, proceed with the update.
		if movie.ID == params["id"] {
			// Remove the existing movie from the slice.
			// Using append to exclude the movie at the current index.
			movies = append(movies[:index], movies[index+1:]...)

			// Declare a new variable to hold the updated movie details.
			var movie Movie

			// Decode the JSON body of the incoming HTTP request into the newly declared movie variable.
			// Notice the use of `&movie`, passing a pointer to the movie variable.
			// This allows json.NewDecoder to populate the movie variable directly.
			_ = json.NewDecoder(r.Body).Decode(&movie)  // Error ignored for brevity, but it's good to handle it.

			// Set the ID of the updated movie to the ID from the URL parameter.
			movie.ID = params["id"]

			// Append the updated movie to the global movies slice.
			movies = append(movies, movie)

			// Encode the updated movie to JSON and write it to the response.
			// This gives the client confirmation of what was updated.
			json.NewEncoder(w).Encode(movie)

			// Return from the function since the update was successful.
			return
		}
	}
}



// Entry point of the program
func main() {
	// Initialize the router
	r := mux.NewRouter()

	// Seed the movies slice with some initial data
	movies = append(movies, Movie{ID: "1", Isbn: "123456", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "123457", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// Registering the handlers with the router
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	//createMovie and updateMovie functions are missing from the original code
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Start the HTTP server
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
