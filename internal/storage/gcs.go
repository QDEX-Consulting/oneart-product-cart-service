package storage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
)

// GCSClient is a helper struct for interacting with Google Cloud Storage.
type GCSClient struct {
	bucketName string
	client     *storage.Client
}

// NewGCSClient initializes a GCS client using environment variables:
//
//   - GCS_BUCKET_NAME: The name of your GCS bucket (e.g., "oneart-product-images")
//   - GOOGLE_APPLICATION_CREDENTIALS: Path to your service account JSON
func NewGCSClient(ctx context.Context) (*GCSClient, error) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME not set")
	}

	// Create the storage client using default credentials.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %v", err)
	}

	return &GCSClient{
		bucketName: bucketName,
		client:     client,
	}, nil
}

// UploadImage uploads the given data (bytes of an image) to GCS under objectName.
// It returns the URL (assuming your bucket or object is publicly readable).
func (g *GCSClient) UploadImage(ctx context.Context, objectName string, data []byte) (string, error) {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)

	// Create a new writer to the object.
	w := obj.NewWriter(ctx)
	defer w.Close()

	// Write the file data.
	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data to GCS: %w", err)
	}

	// Construct a public URL if the bucket or object is publicly accessible.
	// For a private bucket, you'd use signed URLs instead.
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectName)

	return url, nil
}
