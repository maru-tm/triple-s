package storage

import (
	"encoding/xml"

	"triple-s/config"
)

var (
	StorageDir string
	BucketFile string
)

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Message string   `xml:"Message"`
}

// Initialize function to set up storageDir and bucketFile
func InitStorage() {
	StorageDir = config.GetStorageDir()
	BucketFile = StorageDir + "/buckets.csv"
}
