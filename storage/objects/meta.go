package objects

import (
	"encoding/xml"
	"time"
)

// Object represents metadata and details of an object stored in a bucket.
type Object struct {
	BucketName       string    `xml:"bucketName"`
	Key              string    `xml:"key"`
	ContentType      string    `xml:"contentType"`
	Size             string    `xml:"size"`
	LastModifiedTime time.Time `xml:"lastModifiedTime"`
}

type ObjectList struct {
	XMLName xml.Name `xml:"ObjectList"`
	Objects []Object `xml:"Object"`
}
