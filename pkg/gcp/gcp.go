package gcp

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCPUploader struct {
	Bucket    string
	CredsFile string

	GCPClient storage.Client
}

func NewGCPUploader(bucket, credsFile string) GCPUploader {
	gcpClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(credsFile))
	if err != nil {
		panic(err)
	}

	return GCPUploader{
		Bucket:    bucket,
		CredsFile: credsFile,
		GCPClient: *gcpClient,
	}
}

func (g *GCPUploader) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Define your bucket name and object (file) name
	objectName := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)

	// Upload the file to Google Cloud Storage
	bucket := g.GCPClient.Bucket(g.Bucket)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// Make the file publicly accessible (optional)
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	// File uploaded successfully, return its URL
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.Bucket, objectName)
	return fileURL, nil
}
