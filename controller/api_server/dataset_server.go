package main

import (
	"context"
	"dataset/controller"
	log "dataset/logger"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var ctx = context.Background()
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", handler)
	log.Info(ctx, "Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(ctx, "Error starting server: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var start = time.Now()
	request, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, `Error reading request to server`)
		return
	}
	responder(w, request)
	log.Info(context.TODO(), time.Since(start))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var start = time.Now()
	// Parse the multipart form with a max memory of 10 MiB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, "Failed to parse multipart form")
		return
	}
	audioFile, audioHeader, err := r.FormFile("audio")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err, "Invalid audio file")
		return
	}
	defer audioFile.Close()
	// Save the audio file to disk
	audioDst, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), "*"+audioHeader.Filename)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, "Failed to create audio file on server")
		return
	}
	defer audioDst.Close()
	_, err = io.Copy(audioDst, audioFile)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, "Failed to save audio file")
		return
	}
	yamlFile, yamlHeader, err := r.FormFile("yaml")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err, "Invalid YAML file")
		return
	}
	defer yamlFile.Close()
	request, err := io.ReadAll(yamlFile)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, "Unable to read YAML file")
		return
	}
	log.Info(context.TODO(), "Files uploaded successfully:", audioHeader.Filename, yamlHeader.Filename)
	responder(w, request)
	log.Info(context.TODO(), time.Since(start))
}

func responder(w http.ResponseWriter, request []byte) {
	var ctx = context.WithValue(context.Background(), `runType`, `server`)
	var control = controller.NewController(ctx, request)
	var filename, status = control.Process()
	if status.IsErr {
		w.WriteHeader(status.Status)
	} else {
		w.WriteHeader(200)
	}
	var mimeType string
	if strings.HasSuffix(filename, `.csv`) {
		mimeType = "text/csv"
	} else if strings.HasSuffix(filename, `.json`) {
		mimeType = "application/json"
	} else {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	var file *os.File
	file, err := os.Open(filename)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, `File containing results is not found`)
		return
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, `Error writing results to http response`)
		return
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	status := log.Error(context.TODO(), statusCode, err, message)
	http.Error(w, status.String(), statusCode)
}
