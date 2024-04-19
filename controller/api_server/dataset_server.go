package main

import (
	"context"
	"dataset/controller"
	log "dataset/logger"
	"io"
	"net/http"
)

func main() {
	var ctx = context.Background()
	http.HandleFunc("/", echoHandler) // Set the handler function for root path
	log.Info(ctx, "Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil) // Start the server on port 8080
	if err != nil {
		log.Panic(ctx, "Error starting server: ", err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	//_, err := io.Copy(w, r.Body)
	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var control = controller.NewController(content)
	var output = control.Process()
	w.Write(output)
	//if err != nil {
	//	http.Error(w, "Failed to echo request body", http.StatusInternalServerError)
	//}
}
