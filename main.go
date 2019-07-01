package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var gfile string
var result map[string]interface{}
var tempOn bool

func main() {
	port := flag.String("port", "8080", "set the port")
	file := flag.String("file", "file.json", "json file name")
	reRead := flag.Bool("temp", false, "read the json file everytime call endpoint")
	flag.Parse()
	gfile = *file
	tempOn = *reRead

	jsonFile, err := os.Open(gfile)
	if err != nil {
		log.Println(err)
	}
	if jsonFile == nil {
		log.Fatal("there is no file named: ", gfile)
	}
	log.Println("Successfully Opened", gfile)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)

	var router = mux.NewRouter()
	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/", endpoint).Methods("GET")

	log.Println("Running server at", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Pong!")
	json.NewEncoder(w).Encode("pong!")
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

	if tempOn {
		jsonFile, err := os.Open(gfile)
		if err != nil {
			log.Println(err)
		}
		log.Println("Serving...", gfile)
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		json.NewEncoder(w).Encode(result)
		return
	}
	json.NewEncoder(w).Encode(result)
}
