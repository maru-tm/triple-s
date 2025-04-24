package api

import (
	"net/http"
	"strings"

	"triple-s/storage/objects"
	"triple-s/utils"
)

func handlePutObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Split the URL path to get the bucket name and object key
	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid URL. Expected /{BucketName}/{ObjectKey}", http.StatusBadRequest)
		return
	}

	bucketName, objectKey := parts[0], parts[1]
	// Validate the object key
	if !utils.ValidateObjectKey(objectKey) {
		http.Error(w, "400 Bad Request: Invalid object key", http.StatusBadRequest)
		return
	}

	// Create the object using the uploaded file and extracted metadata
	if err := objects.CreateObject(bucketName, objectKey, w, r); err != nil {
		http.Error(w, "500 Internal Server Error: Failed to create object", http.StatusInternalServerError)
		return
	}
}

// handleGetObject retrieves and responds with the object's metadata in XML format.
func handleGetObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid URL. Expected /{BucketName}/{ObjectKey}", http.StatusBadRequest)
		return
	}
	bucketName, objectKey := parts[0], parts[1]

	// Get object metadata
	err := objects.GetObject(bucketName, objectKey, w, r)
	if err != nil {
		http.Error(w, "404 Not Found: Object not found", http.StatusNotFound)
		return
	}

	// Set the response header for XML content
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

// handleDeleteObject removes the object and its metadata.
func handleDeleteObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid URL. Expected /{BucketName}/{ObjectKey}", http.StatusBadRequest)
		return
	}
	bucketName, objectKey := parts[0], parts[1]

	if err := objects.DeleteObject(bucketName, objectKey); err != nil {
		http.Error(w, "500 Internal Server Error: Error deleting object", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
