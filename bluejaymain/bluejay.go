package bluejaymain

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Mark struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
	Created  time.Time `json:"created"`
}
type Marks []Mark

const oneMB int64 = 1048576

// TODO refine this better. Some calls should not kill process
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func returnAllBookmarks(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: returnAllBookmarks")

	marks := getAll()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(marks); err != nil {
		log.Println(err)
	}
}

func addSingleBookrmark(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: addSingleBookrmark")

	var mark Mark
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, oneMB))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &mark); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	m := addMark(mark)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/bookmarks", returnAllBookmarks).Methods("GET")
	router.HandleFunc("/bookmark", addSingleBookrmark).Methods("POST")
	log.Println("Ready to accept connections")
	log.Fatal(http.ListenAndServe("localhost:8081", router))
}

func Main() {
	handleRequests()
	log.Println("Exiting...")
}
