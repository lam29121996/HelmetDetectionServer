package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type config struct {
	Port int `json:"port"`
}

var (
	cfg config

	resCh chan string
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! This is indexHandler of helmet_detection_server!"))
}

func postHelmetDetectionResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	// Let see how camera send the photo to me, file / file path of the photo?
	// no helmet detected -> the event is triggered -> send photo


}

func getHelmetDetectionResultHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		IsHelmetOn bool   `json:"is_helmet_on"`
		PhotoPath  string `json:"photo_path"`
	}

	resp := Response{}

	go func() {
		resCh <- "C:"
	}()

	select {
	case res := <-resCh:
		resp = Response{IsHelmetOn: false, PhotoPath: res}
	case <-time.After(1 * time.Second):
		resp = Response{IsHelmetOn: true, PhotoPath: ""}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}

func init() {
	// Read config.json
	f, err := os.Open("config.json")
	if err != nil {
		log.Panicln("Open config.json failed!", err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panicln("Read config.json failed!", err)
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		log.Panicln("Unmarshal config.json failed!", err)
	}

	resCh = make(chan string)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/helmetDetectionResult", postHelmetDetectionResultHandler).Methods("POST")
	r.HandleFunc("/helmetDetectionResult", getHelmetDetectionResultHandler).Methods("GET")

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
