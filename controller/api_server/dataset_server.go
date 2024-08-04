package main

import (
	"context"
	"dataset/controller"
	log "dataset/logger"
	"io"
	"mime/multipart"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

func main() {
	var ctx = context.Background()
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/request", handler)
	log.Info(ctx, "Server starting on port 7777...")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Panic(ctx, "Error starting server: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var start = time.Now()
	if r.Method != `POST` {
		errorResponse(w, http.StatusMethodNotAllowed, nil, `Only POST method is allowed`)
	}
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
	if r.Method != `POST` {
		errorResponse(w, http.StatusMethodNotAllowed, nil, `Only POST method is allowed`)
	}
	// Parse the multipart form with a max memory of 10 MiB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err, "Failed to parse multipart form")
		return
	}
	// Read form parts, currently tested to handle yaml file and one content file
	var request []byte
	var yamlHeader *multipart.FileHeader
	var dataHeader *multipart.FileHeader
	for key := range r.MultipartForm.File {
		file, header, err2 := r.FormFile(key)
		if err2 != nil {
			errorResponse(w, http.StatusInternalServerError, err2, "Failed to read multipart form")
			return
		}
		defer file.Close()
		if key == "yaml" {
			yamlHeader = header
			request, err = io.ReadAll(file)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err, "Unable to read YAML file")
				return
			}
		} else {
			dataHeader = header
			var target *os.File
			target, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), "*"+header.Filename)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err, "Failed to create audio file on server")
				return
			}
			defer target.Close()
			_, err = io.Copy(target, file)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err, "Failed to save audio file")
				return
			}
		}
	}
	log.Info(context.TODO(), "Files uploaded successfully:", dataHeader.Filename, yamlHeader.Filename)
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
