package objects

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"triple-s/storage"
)

// CreateObjectMeta saves object metadata (content type, size, timestamps).
func CreateObjectMeta(bucketName, key, contentType, size string) error {
	objectFile := filepath.Join(storage.StorageDir, bucketName, "objects.csv")

	// Open the CSV file for reading and writing; create it if it doesn't exist.
	file, err := os.OpenFile(objectFile, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read existing records from the file.
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read records: %w", err)
	}

	// Move the file pointer back to the start to overwrite the content.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	currentTime := time.Now().Format(time.RFC3339)
	modified := false

	// Iterate over the records to update if the key matches.
	for _, record := range records {
		if record[0] == bucketName && record[1] == key {
			// Update the record if it matches the bucket and key.
			record = []string{bucketName, key, contentType, size, currentTime}
			modified = true
		}
		// Write each record back to the file.
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	// If no matching record was found, write a new one.
	if !modified {
		if err := writer.Write([]string{bucketName, key, contentType, size, currentTime}); err != nil {
			return fmt.Errorf("failed to write new record: %w", err)
		}
	}

	// Ensure the file is properly flushed and truncated at the right point.
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

func ListObjects(w http.ResponseWriter, r *http.Request, bucketName string) error {
	objectFile := filepath.Join(storage.StorageDir, bucketName, "/objects.csv")

	// Open the CSV file for reading
	file, err := os.Open(objectFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var objects []Object
	for _, record := range records[1:] { // Assuming the first record is a header
		if record[0] == bucketName {
			lastModifiedTime, _ := time.Parse(time.RFC3339, record[4])
			objects = append(objects, Object{
				BucketName:       record[0],
				Key:              record[1],
				ContentType:      r.Header.Get("Content-Type"),
				Size:             r.Header.Get("Content-Length"),
				LastModifiedTime: lastModifiedTime,
			})
		}
	}
	response := ObjectList{Objects: objects}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
	return nil
}

func DeleteObjectMeta(bucketName, key string) error {
	objectFile := filepath.Join(storage.StorageDir, bucketName, "/objects.csv")

	file, err := os.Open(objectFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	newFile, err := os.Create(objectFile)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	for _, record := range records {
		if record[0] != bucketName || record[1] != key {
			err := writer.Write(record)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
