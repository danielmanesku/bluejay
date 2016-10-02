package bluejaymain

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

type Mark struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
	Created  time.Time `json:"created"`
}
type Marks []Mark

const oneMB int64 = 1048576

func returnAllBookmarks(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: returnAllBookmarks")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	marks := getAll()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(marks); err != nil {
		log.Error(err)
	}
}

func addSingleBookrmark(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: addSingleBookrmark")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var mark Mark
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, oneMB))
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &mark); err != nil {
		log.WithFields(log.Fields{
			"body": body,
		}).Debug("JSON content could not be unmarshaled")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Error(err)
			return
		}
		return
	}

	m := addMark(mark)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Error(err)
		return
	}
}

func deleteSingleBookrmark(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: deleteSingleBookrmark")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("id param must be a number")
		return
	}

	log.Info("DELETE Request for:", id)
	if err := deleteMark(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/bookmarks", returnAllBookmarks).Methods("GET")
	router.HandleFunc("/bookmark", addSingleBookrmark).Methods("POST")
	router.HandleFunc("/bookmark/{id}", deleteSingleBookrmark).Methods("DELETE")
	log.Info("Ready to accept connections")
	log.Fatal(http.ListenAndServe("localhost:8081", router))
}

func Main() {
	handleRequests()
	log.Info("Exiting...")
}
