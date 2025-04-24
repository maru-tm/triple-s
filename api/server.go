package api

import (
	"log"
	"net/http"
	"strings"

	"triple-s/config"
	"triple-s/storage/objects"
)

// StartServer initializes and starts the HTTP server
func StartServer(config *config.Config) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		switch r.Method {
		case http.MethodPut:
			if len(pathParts) == 1 && pathParts[0] != "" {
				handlePutBucket(w, r) // Handle bucket creation
			} else if len(pathParts) == 2 {
				handlePutObject(w, r) // Handle object upload
			} else {
				http.Error(w, "Invalid URL or missing parameters", http.StatusBadRequest)
			}

		case http.MethodGet:

			if r.URL.Path == "/" {
				handleGetBuckets(w, r) // List all buckets
			} else if len(pathParts) == 1 {
				objects.ListObjects(w, r, pathParts[0])
			} else if len(pathParts) == 2 {
				handleGetObject(w, r) // Retrieve a specific object
			} else {
				http.Error(w, "Bucket name or object key is required", http.StatusBadRequest)
			}

		case http.MethodDelete:
			if len(pathParts) == 1 && pathParts[0] != "" {
				handleDeleteBucket(w, r) // Delete a bucket
			} else if len(pathParts) == 2 {
				handleDeleteObject(w, r) // Delete a specific object
			} else {
				http.Error(w, "Invalid URL or missing parameters", http.StatusBadRequest)
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server on the specified port
	log.Printf("Server is listening on port %s", config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatalf("Could not listen on port %s: %v", config.Port, err)
		return err
	}
	return nil
}
