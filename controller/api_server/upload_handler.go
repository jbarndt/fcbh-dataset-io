package main

import (
	"context"
	log "dataset/logger"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

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
	// Append filepath to request
	reqStr := string(request)
	if !strings.HasSuffix(reqStr, "\n") {
		reqStr += "\n"
	}
	reqStr += `uploaded_filepath: ` + audioDst.Name() + "\n"
	request = []byte(reqStr)
	log.Info(context.TODO(), "Files uploaded successfully:", audioHeader.Filename, yamlHeader.Filename)
	responder(w, request)
	log.Info(context.TODO(), time.Since(start))
}
