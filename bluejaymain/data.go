package bluejaymain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

const dataFileName string = "data.json"

func init() {
	// Create empty data file if it does not exist already
	_, err := os.Stat(dataFileName)
	if err != nil {
		if err = ioutil.WriteFile(dataFileName, []byte(`[]`), 0600); err != nil {
			log.WithFields(log.Fields{
				"file name": dataFileName,
			}).Panic("File could not be written. Cannot continue.")
		}
		log.WithFields(log.Fields{
			"file name": dataFileName,
		}).Info("Data file did not exist and it was created")
	}
}

func getAll() Marks {
	var marks Marks
	fileData, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		log.Fatal("Data file is supposed to exist. It should have been created in init func")
	}
	if err := json.Unmarshal(fileData, &marks); err != nil {
		log.Error("Could not unmarshall data file to JSON")
	}
	return marks
}

func addMark(m Mark) Mark {
	calculateNextId := func(marks Marks) int {
		nextId := 0
		for _, m := range marks {
			if m.ID > nextId {
				nextId = m.ID
			}
		}
		return nextId + 1
	}

	allMarks := getAll()
	m.ID = calculateNextId(allMarks)
	m.Created = time.Now().UTC()
	allMarks = append(allMarks, m)
	persistMarks(allMarks)

	return m
}

func deleteMark(id int) error {
	allMarks := getAll()

	removed := false
	// exclude the mark if exists
	for i, m := range allMarks {
		if m.ID == id {
			allMarks = append(allMarks[:i], allMarks[i+1:]...)
			removed = true
		}
	}
	if !removed {
		log.WithFields(log.Fields{
			"id": id,
		}).Info("Bookrmark not found")
		return errors.New("not found")
	}
	persistMarks(allMarks)

	return nil
}

func persistMarks(marks Marks) {
	jsonData, err := json.Marshal(marks)
	if err != nil {
		log.WithFields(log.Fields{
			"field": marks,
		}).Error("Could not marshal data structure to JSON")
	}
	if err := ioutil.WriteFile(dataFileName, jsonData, 0600); err != nil {
		log.Error("Could not persist data.")
	}
}
