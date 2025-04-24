package buckets

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"triple-s/storage"
)

// CreateRootDirectory creates the {data} directory and the buckets.csv file only if they don't exist
func CreateRootDirectory() error {
	rootDir := filepath.Join(storage.StorageDir)

	// Check if the {data} directory exists
	dirInfo, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(rootDir, 0o755)
		if err != nil {
			return err
		}
	} else if !dirInfo.IsDir() {
		// If it's not a directory, return an error
		return fmt.Errorf("%s exists but is not a directory", rootDir)
	}

	// Check if the buckets.csv file exists
	objectsFile := filepath.Join(rootDir, "buckets.csv")
	if _, err := os.Stat(objectsFile); os.IsNotExist(err) {
		// Create the buckets.csv file
		file, err := os.Create(objectsFile)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write the header row to the CSV file
		writer := csv.NewWriter(file)
		defer writer.Flush()
		headers := []string{"BucketName", "CreationTime", "LastModifiedTime", "Status"}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	return nil
}

// Create the bucket directory
func CreateBucketDirectory(bucketName string) error {
	// storageDir := config.GetStorageDir()
	bucketDir := filepath.Join(storage.StorageDir, bucketName)

	// Create the bucket directory
	err := os.MkdirAll(bucketDir, 0o755)
	if err != nil {
		return err
	}

	// Create the objects.csv file
	objectsFile := filepath.Join(bucketDir, "objects.csv")
	file, err := os.Create(objectsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the header row to the CSV file
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{"BucketName", "ObjectKey", "ContentType", "Size", "LastModifiedTime"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Store object metadata
	err = CreateBucketMeta(bucketName)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBucketDirectory removes the specified bucket directory if it contains only 'objects.csv'.
func DeleteBucketDirectory(bucketName string) error {
	// Check if the bucket exists
	exists, err := storage.BucketExists(bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	bucketDir := filepath.Join(storage.StorageDir, bucketName)

	// Check if the bucket directory exists
	_, err = os.Stat(bucketDir)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	// Check if the bucket is empty (only contains 'objects.csv')
	if storage.IsBucketEmpty(bucketName) {
		err = DeleteBucketMeta(bucketName)
		if err != nil {
			return err
		}
		return os.RemoveAll(bucketDir)
	}

	// If the directory has more than 'objects.csv', return an error
	return os.ErrExist
}
