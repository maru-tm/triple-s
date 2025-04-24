package buckets

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"triple-s/storage"
)

func CreateBucketMeta(name string) error {
	// Open the CSV file for appending
	file, err := os.OpenFile(storage.BucketFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get the current time for creation and modification
	currentTime := time.Now().Format(time.RFC3339)
	status := "Available"

	// Write the bucket name and timestamps to the CSV
	if err = writer.Write([]string{name, currentTime, currentTime, status}); err != nil {
		return err
	}

	return nil
}

func ListBuckets() ([]Bucket, error) {
	// Открыть CSV файл для чтения
	file, err := os.Open(storage.BucketFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, record := range records[1:] {
		creationTime, _ := time.Parse(time.RFC3339, record[1])
		lastModifiedTime, _ := time.Parse(time.RFC3339, record[2])
		buckets = append(buckets, Bucket{
			Name:             record[0],
			CreationTime:     creationTime,
			LastModifiedTime: lastModifiedTime,
			Status:           record[3],
		})
	}
	return buckets, nil
}

func DeleteBucketMeta(name string) error {
	file, err := os.Open(storage.BucketFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	newFile, err := os.Create(storage.BucketFile)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	for _, record := range records {
		// Skip writing the record if it matches the bucket to be deleted
		if record[0] == name {
			continue
		}

		// Write the remaining records to the new file
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
