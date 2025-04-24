package storage

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrBucketExists = errors.New("bucket already exists")
	ErrObjectExists = errors.New("object already exists")
)

// BucketExists checks if a bucket with the given name exists.
func BucketExists(name string) (bool, error) {
	// Open the bucket file for reading
	file, err := os.Open(BucketFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read all records from the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	// Check if the specified bucket name exists in the records
	for _, record := range records {
		if record[0] == name && record[3] != "Deleted" {
			return true, nil
		}
	}
	return false, nil
}

// IsBucketEmpty checks if the bucket has any objects other than 'objects.csv'.
func IsBucketEmpty(bucketName string) bool {
	bucketDir := filepath.Join(StorageDir, bucketName)

	files, err := os.ReadDir(bucketDir)
	if err != nil {
		return false
	}

	if len(files) == 0 {
		return true
	}

	if len(files) == 1 && files[0].Name() == "objects.csv" {
		return true
	}

	return false
}

// ObjectExists checks if an object with the given key in the specified bucket exists.
func ObjectExists(bucketName, key string) (bool, error) {
	objectFile := filepath.Join(StorageDir, bucketName, "objects.csv")

	// Open the object file for reading
	file, err := os.Open(objectFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read all records from the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	// Check if the specified object exists in the records
	for _, record := range records {
		if record[0] == bucketName && record[1] == key {
			return true, nil
		}
	}
	return false, nil
}
