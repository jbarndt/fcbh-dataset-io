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
		errorResponse(w, err, `Error reading request to server`)
		return
	}
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
	file, err = os.Open(filename)
	if err != nil {
		errorResponse(w, err, `File containing results is not found`)
		return
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		errorResponse(w, err, `Error writing results to http response`)
		return
	}
	log.Info(context.TODO(), time.Since(start))
}

func errorResponse(w http.ResponseWriter, err error, message string) {
	status := log.Error(context.TODO(), 500, err, message)
	_, err2 := w.Write([]byte(status.String()))
	if err2 != nil {
		log.Fatal(context.TODO(), err2, message)
	}
}

/*
func main() {
    // Set up a route and attach handler
    http.HandleFunc("/download", fileDownloadHandler)

    // Start the server
    log.Println("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
    // Specify the path to the file
    filePath := "path/to/your/file.txt"

    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        // If the file does not open, send an appropriate response
        http.Error(w, "File not found.", 404)
        return
    }
    defer file.Close()

    // Set the header to ensure the downloaded file has the correct name
    w.Header().Set("Content-Disposition", "attachment; filename=\"" + file.Name() + "\"")
    w.Header().Set("Content-Type", "application/octet-stream")

    // Copy the file to the response writer
    if _, err := io.Copy(w, file); err != nil {
        log.Println("Error writing file to response:", err)
    }
}
*/
