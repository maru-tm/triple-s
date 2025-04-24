package objects

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"triple-s/storage"
)

func CreateObject(bucketName, objectKey string, w http.ResponseWriter, r *http.Request) error {
	// Construct the object path
	objectPath := filepath.Join(storage.StorageDir, bucketName, objectKey)
	// Open file
	objectFile, err := os.Create(objectPath)
	if err != nil {
		return err
	}
	defer objectFile.Close()

	// Copy the content of the Request body to the object
	io.Copy(objectFile, r.Body)

	// Extract the size and the content type from the header
	size := r.Header.Get("Content-Length")
	typeMime := r.Header.Get("Content-Type")

	// Store object metadata
	if err := CreateObjectMeta(bucketName, objectKey, typeMime, size); err != nil {
		return fmt.Errorf("failed to store object metadata: %w", err)
	}

	// Write response
	response := Object{
		BucketName:       bucketName,
		Key:              objectKey,
		ContentType:      typeMime,
		Size:             size,
		LastModifiedTime: time.Now(),
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
	return nil
}

// DeleteObject removes the object and updates the metadata.
func DeleteObject(bucketName, objectKey string) error {
	objectPath := filepath.Join(storage.StorageDir, bucketName, objectKey)

	if err := os.Remove(objectPath); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return DeleteObjectMeta(bucketName, objectKey)
}

// GetObject copies the content of the object to the ResponseWriter.
func GetObject(bucketName, objectKey string, w http.ResponseWriter, r *http.Request) error {
	objectPath := filepath.Join(storage.StorageDir, bucketName, objectKey)
	objectFile, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, objectFile)
	if err != nil {
		return err
	}

	size := r.Header.Get("Content-Length")
	typeMime := r.Header.Get("Content-Type")

	response := Object{
		BucketName:       bucketName,
		Key:              objectKey,
		ContentType:      typeMime,
		Size:             size,
		LastModifiedTime: time.Now(),
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
	return nil
}
