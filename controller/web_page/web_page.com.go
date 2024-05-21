package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func main() {
	http.HandleFunc("/upload", handleUpload)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the upload directory if it doesn't exist
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the uploaded file to disk
	filePath := filepath.Join(uploadDir, handler.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the uploaded file
	processedFilePath := "processed_" + handler.Filename

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename="+processedFilePath)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Open the processed file for reading
	processedFile, err := os.Open(processedFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer processedFile.Close()

	// Stream the processed file to the client
	_, err = io.Copy(w, processedFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*

Here's how the code works:

1. The `handleUpload` function is registered as the handler for the `/upload` route.
2. When a POST request is received at `/upload`, the `handleUpload` function is called.
3. The `r.ParseMultipartForm` function is used to parse the multipart form data, which includes the uploaded file.
4. The `r.FormFile` function retrieves the uploaded file from the form data.
5. The uploaded file is saved to the `./uploads` directory on the server using `os.OpenFile` and `io.Copy`.
6. In this example, we simulate file processing by creating a new file with the prefix `"processed_"` prepended to the original filename.
7. The appropriate headers are set using `w.Header().Set` to indicate that the response body contains a downloadable file.
8. The processed file is opened for reading using `os.Open`.
9. The processed file is streamed to the client using `io.Copy(w, processedFile)`.

Note that in this example, we're not performing any actual file processing. You would need to replace the line `processedFilePath := "processed_" + handler.Filename` with your own file processing logic.

Also, this example doesn't include error handling or cleanup for the uploaded and processed files. In a production environment, you should handle errors gracefully and clean up any temporary files or directories.​​​​​​​​​​​​​​​​
*/
