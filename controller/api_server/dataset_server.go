package main

import (
	"context"
	"dataset/controller"
	log "dataset/logger"
	"io"
	"net/http"
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
	content, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err, `Error reading request to server`)
	} else {
		var control = controller.NewController(content)
		var output = control.Process()
		_, err = w.Write(output)
		if err != nil {
			errorResponse(w, err, `Error writing server response`)
		}
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
