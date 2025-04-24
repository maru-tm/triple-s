package buckets

import (
	"encoding/xml"
	"time"
)

type Bucket struct {
	Name             string    `xml:"Name"`
	CreationTime     time.Time `xml:"CreationTime"`
	LastModifiedTime time.Time `xml:"LastModifiedTime"`
	Status           string    `xml:"Status"`
}

type BucketList struct {
	XMLName xml.Name `xml:"BucketList"`
	Buckets []Bucket `xml:"Bucket"`
}
