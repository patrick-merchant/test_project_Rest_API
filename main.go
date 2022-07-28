package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Mimics the structure of the json placeholder endpoint
type Post struct {
	UserId string `json:"UserId"`
	Id     string `json:"Id"`
	Title  string `json:"Title"`
	Body   string `json:"Body"`
}

var Posts []Post

var visitCounter = 0

// var postCount = len(Posts)

// Visit tracker - closure functions & the wrapper technique.
func trackVisits(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		visitCounter++
		res := "no visits"
		if visitCounter == 1 {
			res = fmt.Sprint("1 visit !")
		} else {
			res = fmt.Sprintf("%v visits !", visitCounter)
		}

		fmt.Println(res)

		handler(w, r);
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

	var myRouter = mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", trackVisits(homePageHandler))
	myRouter.HandleFunc("/posts", trackVisits(returnAllPosts))
	myRouter.HandleFunc("/post", createNewPost).Methods("POST")
	myRouter.HandleFunc("/post/{id}", deletePost).Methods("DELETE")
	myRouter.HandleFunc("/post/{id}", updatePost).Methods("PUT")
	myRouter.HandleFunc("/post/{id}", trackVisits(returnSinglePost))

	log.Fatal(http.ListenAndServe(":8083", myRouter))
}

// Retrieving all posts
func returnAllPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPosts")
	json.NewEncoder(w).Encode(Posts)
}

// Retrieving a single post
func returnSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, post := range Posts {
		if post.Id == key {
			json.NewEncoder(w).Encode(post)
		}
	}
}

// CREATE
func createNewPost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(reqBody, &post)
	Posts = append(Posts, post)

	json.NewEncoder(w).Encode(post)
}

// DELETE
func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	for index, post := range Posts {
		if post.Id == id {
			Posts = append(Posts[:index], Posts[index+1:]...)
		}
	}
}

// UPDATE
func updatePost(w http.ResponseWriter, r *http.Request) {

	// Parses HTTP request body and stores json in post variable (like CREATE)
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(reqBody, &post)
	json.NewEncoder(w).Encode(post)

	// Loops over posts in Posts array and updates the selected post
	vars := mux.Vars(r)
	id := vars["id"]
	for index, value := range Posts {
		if value.Id == id {
			Posts[index] = post
			fmt.Println(Posts)
		}
	}
}

func main() {

	fmt.Println("Rest API v2.0 - Mux Routers")
	Posts = []Post{
		Post{UserId: "1", Id: "1", Title: "A Post", Body: "This is a post"},
		Post{UserId: "2", Id: "2", Title: "Another Post", Body: "This is another post"},
	}
	handleRequests()
}
