package bluejaymain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const dataFileName string = "data.json"

func init() {
	// Create empty data file if it does not exist already
	_, err := os.Stat(dataFileName)
	if err != nil {
		if err = ioutil.WriteFile(dataFileName, []byte(`[]`), 0600); err != nil {
			log.Fatal("File could not be written. Cannot continue.")
		}
	}
	log.Printf("%s did not exist and it was created", dataFileName)
}

func getAll() Marks {
	var marks Marks
	fileData, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		// return make(Marks, 0), err
		log.Fatal("Data file is supposed to exist. It should have been created in init func")
	}
	if err := json.Unmarshal(fileData, &marks); err != nil {
		log.Fatal("Error: could not unmarshall data file to JSON") //TODO investigate proper logging
		// return nil, err
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
	allMarks = append(allMarks, m)

	// persist data
	jsonData, err := json.Marshal(allMarks)              //TODO check err with log level
	err = ioutil.WriteFile(dataFileName, jsonData, 0600) //TODO should be err :=
	check(err)

	return m
}
