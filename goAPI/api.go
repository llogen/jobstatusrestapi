package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Define PORT and FILENAME
const (
	PORT     = ":3005"
	FILENAME = "go_jobstatus.json"
)

// DB data structure
type JsonData struct {
	Entrys []Entry
}

// Entry structure
type Entry struct {
	ID     int
	Status bool
}

// Function that reads the status out of the DB for the requested job
func readjobstatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("readjobstatus was triggered")
	// retrive the parameter id and convert it to int
	vars := mux.Vars(r)
	var reqID = vars["id"]
	reqid, err := strconv.Atoi(reqID)
	if err != nil {
		panic("Could not convert string to int.")
	}
	// check the DB for existing entry and status if exists
	exists, status := readDBdata(reqid, FILENAME)
	if exists == "" {
		// if not successful
		json.NewEncoder(w).Encode("Job does not exist!")
	} else {
		// if successful
		json.NewEncoder(w).Encode(status)
	}
}

// Function that updates a job status in the DB
func updatejobstatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updatejobstatus was triggered")
	// decode the json body to access the body data
	jsonbody := json.NewDecoder(r.Body)
	var entry Entry
	err := jsonbody.Decode(&entry)
	if err != nil {
		panic("Could not decode the JSON Body.")
	}

	// unmarshal the DB json file to update the status
	var jsonData JsonData
	jsonFile, err := ioutil.ReadFile(FILENAME)
	if err != nil {
		panic("Could not open the DB file.")
	}
	json.Unmarshal(jsonFile, &jsonData)

	// if entry exists change it and write it back to the DB json file
	var found = 0
	for i := 1; i < len(jsonData.Entrys)+1; i++ {
		if jsonData.Entrys[i-1].ID == entry.ID {
			found = i
		}
	}
	if found > 0 {
		jsonData.Entrys[found-1].Status = entry.Status
		errmsg := marshalAndWrite(jsonData, FILENAME)
		if errmsg != "" {
			// if not successful
			json.NewEncoder(w).Encode(errmsg)
			return
		}
		// if successful
		json.NewEncoder(w).Encode("Status edited!")
	} else {
		// if not successful
		json.NewEncoder(w).Encode("Job does not exist!")
	}

}

// Function that adds a job to the DB
func addjobstatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addjobstatus was triggered")
	// decode the json body to access the body data
	jsonbody := json.NewDecoder(r.Body)
	var entry Entry
	err := jsonbody.Decode(&entry)
	if err != nil {
		panic("Could not decode the JSON Body.")
	}

	// check the DB for existing entry and status if exists
	exists, _ := readDBdata(entry.ID, FILENAME)
	if exists == "" {
		// unmarshal the DB json file to update the status
		var jsonData JsonData
		jsonFile, err := ioutil.ReadFile(FILENAME)
		if err != nil {
			panic("Could not open the DB file.")
		}
		json.Unmarshal(jsonFile, &jsonData)
		// Create data that should be added, append it and write it back to the DB json file
		var appendData = &Entry{ID: entry.ID, Status: entry.Status}
		jsonData.Entrys = append(jsonData.Entrys, *appendData)
		errmsg := marshalAndWrite(jsonData, FILENAME)
		if errmsg != "" {
			json.NewEncoder(w).Encode(errmsg)
			return
		}
		// if successful
		json.NewEncoder(w).Encode("Jobstatus added!")
	} else {
		// if not successful
		json.NewEncoder(w).Encode("The Job already exists!")
	}
}

// Function that removes a job out of the DB
func removejobstatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("removejobstatus was triggered")
	// retrive the parameter id and convert it to int
	vars := mux.Vars(r)
	var reqID = vars["id"]
	reqid, err := strconv.Atoi(reqID)
	if err != nil {
		panic("Could not convert string to int.")
	}

	// check the DB for existing entry and status if exists
	exists, _ := readDBdata(reqid, FILENAME)
	if exists == "" {
		json.NewEncoder(w).Encode("Job does not exist!")
	} else {
		// unmarshal the DB json file to update the status
		var jsonData JsonData
		jsonFile, err := ioutil.ReadFile(FILENAME)
		if err != nil {
			panic("Could not open the DB file.")
		}
		json.Unmarshal(jsonFile, &jsonData)
		// search for the correct job, remove the whole slice and write it back to the DB json file
		for i := 0; i < len(jsonData.Entrys); i++ {
			if jsonData.Entrys[i].ID == reqid {
				jsonData.Entrys = append(jsonData.Entrys[:i], jsonData.Entrys[i+1:]...)
				errmsg := marshalAndWrite(jsonData, FILENAME)
				if errmsg != "" {
					// if not successful
					json.NewEncoder(w).Encode(errmsg)
					return
				}
				// if successful
				json.NewEncoder(w).Encode("Job removed!")
			}
		}
	}
}

func main() {
	fmt.Println("Starting RestAPI and listening on Port", PORT)
	handleRequests()
}

func handleRequests() {
	// new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/")
	myRouter.HandleFunc("/readjobstatus/{id}", readjobstatus)
	myRouter.HandleFunc("/updatejobstatus/", updatejobstatus).Methods("POST")
	myRouter.HandleFunc("/addjobstatus/", addjobstatus).Methods("POST")
	myRouter.HandleFunc("/removejobstatus/{id}", removejobstatus).Methods("POST")

	log.Fatal(http.ListenAndServe(PORT, myRouter))
}

// Function that reads the DB json file and searches for the job with the JobID ID
func readDBdata(ID int, filename string) (string, bool) {
	// Read the json file and unmarshal it
	var jsonData JsonData
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Could not open the DB file.")
	}
	json.Unmarshal(jsonFile, &jsonData)
	// Go through the json data and check against the ID
	// Return a non empty string and status if ID exists
	for i := 0; i < len(jsonData.Entrys); i++ {
		if jsonData.Entrys[i].ID == ID {
			return "exists", jsonData.Entrys[i].Status
		}
	}
	// Return an empty string if the ID not exists
	return "", false
}

// Function that marshal a given data and write it back to the DB json file
func marshalAndWrite(jsonData JsonData, filename string) string {
	// Marshal the data
	newjsonData, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return "Could not marshal the json data."
	}
	// Write the marshaled data to file
	err = ioutil.WriteFile(filename, newjsonData, 0644)
	if err != nil {
		return "Could not write new data into the DB file."
	}
	// Return empty string if everything worked fine
	return ""
}
