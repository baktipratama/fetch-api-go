package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Response struct {
	Data int64 `json:"data"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Review_data struct {
	Reviewid  string `json:"reviewid"`
	Title string `json:"title"`
	Url  string `json:"url"`
	Score string `json:"score"`
	Artists string `json:"artists"`
	Genres string `json:"genres"`
	Labels string `json:"labels"`
	Pub_year string `json:"pub_year"`
	Content  string `json:"content"`
}

func main() {
	db, err = gorm.Open(sqlite.Open("review.db"), &gorm.Config{})

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/test", testController).Methods("GET")
	myRouter.HandleFunc("/api/reviews", getAllReviews).Methods("POST")
	myRouter.HandleFunc("/api/reviews_by_score", getReviewByScore).Methods("POST")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func testController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fib := fibbonaci(100)
	resp := Response{Data: fib}
	json.NewEncoder(w).Encode(resp)
}

func fibbonaci(n int) int64 {
	f := make([]int64, n+1, n+2)
	if n < 2 {
		f = f[0:2]
	}
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}

func getAllReviews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get all reviews")
	limit := r.FormValue("limit")
	limit_int, err := strconv.Atoi(limit)

	reviews := []Review_data{}
	var query = db.Table("reviews")
	query.Select("reviews.reviewid, reviews.title, reviews.url, reviews.score, (SELECT GROUP_CONCAT(artists.artist)) as artists, (SELECT GROUP_CONCAT(genres.genre)) as genres, (SELECT GROUP_CONCAT(labels.label)) as labels, reviews.pub_year, content.content")
	query.Joins("join labels on labels.reviewid = reviews.reviewid")
	query.Joins("join artists on artists.reviewid = reviews.reviewid")
	query.Joins("join content on content.reviewid = reviews.reviewid")
	query.Joins("join genres on genres.reviewid = reviews.reviewid")
	query.Group("reviews.reviewid")
	query.Limit(limit_int)
	query.Scan(&reviews)

	res := Result{Code: 200, Data: reviews, Message: "Success get reviews"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getReviewByScore(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get all reviews")
	score := r.FormValue("score")

	reviews := []Review_data{}
	var query = db.Table("reviews")
	query.Select("reviews.reviewid, reviews.title, reviews.url, reviews.score, (SELECT GROUP_CONCAT(artists.artist)) as artists, (SELECT GROUP_CONCAT(genres.genre)) as genres, (SELECT GROUP_CONCAT(labels.label)) as labels, reviews.pub_year, content.content")
	query.Joins("join labels on labels.reviewid = reviews.reviewid")
	query.Joins("join artists on artists.reviewid = reviews.reviewid")
	query.Joins("join content on content.reviewid = reviews.reviewid")
	query.Joins("join genres on genres.reviewid = reviews.reviewid")
	query.Where("score >= ?", score)
	query.Group("reviews.reviewid")
	query.Scan(&reviews)

	res := Result{Code: 200, Data: reviews, Message: "Success get reviews"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}
