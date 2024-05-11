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

func responder(w http.ResponseWriter, request []byte) {
	var control = controller.NewController(request)
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
