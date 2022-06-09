package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type config struct {
	DebugMode                int    `json:"debug_mode"`
	Port                     int    `json:"port"`
	Timeoutms                int    `json:"timeout(ms)"`
	ImagesHiuMingFolderPath  string `json:"images_hiu_ming_folder_path"`
	ImagesHiuKwongFolderPath string `json:"images_hiu_kwong_folder_path"`
	RecordFrom               string `json:"record_from"`
	recordFromHour           int
	recordFromMinute         int
	RecordTo                 string `json:"record_to"`
	recordToHour             int
	recordToMinute           int
	CaptureIntervalms        int `json:"capture_interval(ms)"`
}

var (
	cfg config

	imagesHiuMingLastUpdateAt  time.Time
	imagesHiuKwongLastUpdateAt time.Time
)

func toFileName(t time.Time, from string) string {
	str := t.Format(time.RFC3339)
	str = str[:19]
	str = strings.ReplaceAll(str, "T", "_")
	str = strings.ReplaceAll(str, ":", "_")

	if from == "hiuMing" {
		path := fmt.Sprintf("%s/%s", cfg.ImagesHiuMingFolderPath, str[:10])
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		return fmt.Sprintf("%s/%s/%s.jpg", cfg.ImagesHiuMingFolderPath, str[:10], str)
	}

	if from == "hiuKwong" {
		path := fmt.Sprintf("%s/%s", cfg.ImagesHiuKwongFolderPath, str[:10])
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		return fmt.Sprintf("%s/%s/%s.jpg", cfg.ImagesHiuKwongFolderPath, str[:10], str)
	}

	return ""
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! This is indexHandler of helmet_detection_server!"))
}

func createImage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if now.Hour() < cfg.recordFromHour || (now.Hour() == cfg.recordFromHour && now.Minute() < cfg.recordFromMinute) {
		log.Println("too early!")
		return
	}
	if (now.Hour() == cfg.recordToHour && now.Minute() > cfg.recordToMinute) || now.Hour() > cfg.recordToHour {
		log.Println("too late!")
		return
	}

	from := r.URL.Query().Get("from")

	switch from {
	case "hiuMing":
		if now.Sub(imagesHiuMingLastUpdateAt) <= time.Duration(cfg.CaptureIntervalms)*time.Millisecond {
			return
		}
		imagesHiuMingLastUpdateAt = now
	case "hiuKwong":
		if now.Sub(imagesHiuKwongLastUpdateAt) <= time.Duration(cfg.CaptureIntervalms)*time.Millisecond {
			return
		}
		imagesHiuKwongLastUpdateAt = now
	default:
		log.Println("without from!")
	}

	file, err := os.Create(toFileName(now, from))
	if err != nil {
		log.Println("os.Create failed!", err)
	}

	log.Printf("Images from: %s created!\n", from)

	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println("io.Copy failed!", err)
	}
}

func helmetDetectionResult(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()
	from := r.URL.Query().Get("from")

	if from == "" {
		w.Write([]byte("from required!"))
		return
	}

	type Response struct {
		IsHelmetOn bool   `json:"is_helmet_on"`
		ImagePath  string `json:"image_path"`
	}

	resp := Response{}

	time.Sleep(time.Duration(cfg.Timeoutms) * time.Millisecond)

	switch from {
	case "hiuMing":
		if time.Since(imagesHiuMingLastUpdateAt) <= time.Duration(cfg.Timeoutms)*time.Millisecond {
			resp = Response{IsHelmetOn: false, ImagePath: toFileName(imagesHiuMingLastUpdateAt, from)}
		} else {
			resp = Response{IsHelmetOn: true, ImagePath: ""}
		}
	case "hiuKwong":
		if time.Since(imagesHiuKwongLastUpdateAt) <= time.Duration(cfg.Timeoutms)*time.Millisecond {
			resp = Response{IsHelmetOn: false, ImagePath: toFileName(imagesHiuKwongLastUpdateAt, from)}
		} else {
			resp = Response{IsHelmetOn: true, ImagePath: ""}
		}
	default:
		w.Write([]byte("invalid from!"))
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("GET request handled sucessfully in %s.", time.Since(startAt))
}

func debug() {
	call := func(urlPath, method string) error {
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		img, err := os.Open("Order. W870672511 Cancelled.png")
		if err != nil {
			log.Println(err)
		}

		defer img.Close()

		req, err := http.NewRequest(method, urlPath, img)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "image")

		rsp, _ := client.Do(req)
		if rsp.StatusCode != http.StatusOK {
			log.Printf("Request failed with response code: %d", rsp.StatusCode)
		}

		return nil

	}

	reader := bufio.NewReader(os.Stdin)

	for {
		reader.ReadString('\n')

		for i := 0; i < 3; i++ {
			err := call("http://localhost:8080/createImage?from=hiuMing", "POST")
			log.Println("client called from hiuMing!")
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}

		err := call("http://localhost:8080/createImage?from=hiuKwong", "POST")
		log.Println("client called from hiuKwong!")
		if err != nil {
			log.Println(err)
		}
	}
}

func init() {
	// Read config.json
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panicln("Read config.json failed!", err)
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		log.Panicln("Unmarshal config.json failed!", err)
	}

	strs := strings.Split(cfg.RecordFrom, ":")

	cfg.recordFromHour, err = strconv.Atoi(strs[0])
	if err != nil {
		log.Panicln("strconv.Atoi failed!", err)
	}
	cfg.recordFromMinute, err = strconv.Atoi(strs[1])
	if err != nil {
		log.Panicln("strconv.Atoi failed!", err)
	}

	strs = strings.Split(cfg.RecordTo, ":")

	cfg.recordToHour, err = strconv.Atoi(strs[0])
	if err != nil {
		log.Panicln("strconv.Atoi failed!", err)
	}
	cfg.recordToMinute, err = strconv.Atoi(strs[1])
	if err != nil {
		log.Panicln("strconv.Atoi failed!", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/createImage", createImage).Methods("POST")
	r.HandleFunc("/helmetDetectionResult", helmetDetectionResult).Methods("GET")
	r.PathPrefix("/hiuMingImages").Handler(http.StripPrefix("/hiuMingImages", http.FileServer(http.Dir(cfg.ImagesHiuMingFolderPath)))).Methods("GET")
	r.PathPrefix("/hiuKwongImages").Handler(http.StripPrefix("/hiuKwongImages", http.FileServer(http.Dir(cfg.ImagesHiuKwongFolderPath)))).Methods("GET")

	if cfg.DebugMode != 0 {
		go debug()
	}

	log.Printf("This server only save no helmet photo from %s to %s.\n", cfg.RecordFrom, cfg.RecordTo)
	log.Printf("Listen & Serve at port: %d", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
