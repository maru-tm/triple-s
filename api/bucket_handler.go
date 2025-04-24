package api

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"triple-s/storage"
	"triple-s/storage/buckets"
	"triple-s/utils"
)

// handlePutBucket creates a new bucket if it doesn't already exist.
func handlePutBucket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[1:]
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	if !utils.ValidateBucketName(bucketName) {
		http.Error(w, "400 Bad Request: Invalid bucket name", http.StatusBadRequest)
		return
	}

	exists, err := storage.BucketExists(bucketName)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Error checking bucket existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "409 Conflict: Bucket already exists", http.StatusConflict)
		return
	}

	if err := buckets.CreateBucketDirectory(bucketName); err != nil {
		http.Error(w, "500 Internal Server Error: Error creating bucket", http.StatusInternalServerError)
		return
	}

	response := buckets.Bucket{Name: bucketName, Status: "Created"}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
}

// handleGetBuckets lists all available buckets.
func handleGetBuckets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketData, err := buckets.ListBuckets()
	if err != nil {
		fmt.Printf("Error retrieving bucket list: %v\n", err)
		http.Error(w, "500 Internal Server Error: Error retrieving bucket list", http.StatusInternalServerError)
		return
	}

	var bucketList buckets.BucketList
	for _, data := range bucketData {
		bucketList.Buckets = append(bucketList.Buckets, buckets.Bucket{
			Name:             data.Name,
			CreationTime:     data.CreationTime,
			LastModifiedTime: data.LastModifiedTime,
		})
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(bucketList); err != nil {
		http.Error(w, "500 Internal Server Error: Unable to encode response", http.StatusInternalServerError)
	}
}

// handleDeleteBucket removes a bucket if it exists and is empty.
func handleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[1:]
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	exists, err := storage.BucketExists(bucketName)
	if err != nil {
		fmt.Printf("Error checking bucket existence: %v\n", err)
		http.Error(w, "500 Internal Server Error: Error checking bucket existence", http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		xml.NewEncoder(w).Encode(storage.ErrorResponse{Message: "Bucket not found"})
		return
	}

	if !storage.IsBucketEmpty(bucketName) {
		w.WriteHeader(http.StatusConflict)
		xml.NewEncoder(w).Encode(storage.ErrorResponse{Message: "Bucket is not empty"})
		return
	}

	if err := buckets.DeleteBucketDirectory(bucketName); err != nil {
		fmt.Printf("Error deleting bucket: %v\n", err)
		http.Error(w, "500 Internal Server Error: Error deleting bucket", http.StatusInternalServerError)
		return
	}

	response := buckets.Bucket{Name: bucketName, Status: "Deleted"}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusNoContent)
	xml.NewEncoder(w).Encode(response)
}
