package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type config struct {
	Port           int    `json:"port"`
	Timeoutms      int    `json:"timeout(ms)"`
	ImagesFilePath string `json:"images_file_path"`
}

var (
	cfg config

	resCh chan string
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! This is indexHandler of helmet_detection_server!"))
}

func postHelmetDetectionResultHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL)

	// Let see how camera send the photo to me, file / file path of the photo?
	// no helmet detected -> the event is triggered -> send photo

	resCh <- r.RequestURI
}

func createImage(w http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the photo key - First Approach
	file, h, err := request.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpfile, err := os.Create("./images/" + h.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func helmetDetectionResult(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()

	type Response struct {
		IsHelmetOn bool   `json:"is_helmet_on"`
		PhotoPath  string `json:"photo_path"`
	}

	resp := Response{}

	select {
	case res := <-resCh:
		resp = Response{IsHelmetOn: false, PhotoPath: res}
	case <-time.After(time.Duration(cfg.Timeoutms) * time.Millisecond):
		resp = Response{IsHelmetOn: true, PhotoPath: ""}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("GET request handled sucessfully in %s.", time.Since(startAt))
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

	// Make resCh
	resCh = make(chan string)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	// r.HandleFunc("/helmetDetectionResult", postHelmetDetectionResultHandler).Methods("POST")
	r.HandleFunc("/createImage", createImage).Methods("POST")
	r.HandleFunc("/helmetDetectionResult", helmetDetectionResult).Methods("GET")
	r.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir(cfg.ImagesFilePath)))).Methods("GET")

	// Client routine
	go func() {
		time.Sleep(3 * time.Second)

		call := func(urlPath, method string) error {
			client := &http.Client{
				Timeout: time.Second * 10,
			}

			// New multipart writer.
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			fw, err := writer.CreateFormFile("photo", "five000000.png")
			if err != nil {
				return err
			}

			file, err := os.Open("five000000.png")
			if err != nil {
				return err
			}

			_, err = io.Copy(fw, file)
			if err != nil {
				return err
			}

			writer.Close()

			req, err := http.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
			if err != nil {
				return err
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			rsp, _ := client.Do(req)
			if rsp.StatusCode != http.StatusOK {
				log.Printf("Request failed with response code: %d", rsp.StatusCode)
			}

			return nil
		}

		err := call("http://localhost:8080/createImage", "POST")
		if err != nil {
			log.Println(err)
		}
		log.Println("called!")
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
