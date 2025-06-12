package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

type ServiceAccountKeyJSON struct {
	Type         string `json:"type"`
	ProjectID    string `json:"type"`
	PrivateKeyID string `json:"type"`
	PrivateKey   string `json:"type"`
	ClientEmail  string `json:"type"`
	ClientID     string `json:"type"`
	AuthURI      string `json:"type"`
	TokenURI     string `json:"type"`
	AuthProvider string `json:"type"`
}

type GCSClient struct {
	ServiceAccountKeyJSON ServiceAccountKeyJSON
	BucketName            string
}

type IGCSClient interface {
	UploadFile(context.Context, string, []byte) (string, error)
}

func NewGCSClient(serviceAccountKeyJSON ServiceAccountKeyJSON, bucketName string) IGCSClient {
	return &GCSClient{
		ServiceAccountKeyJSON: serviceAccountKeyJSON,
		BucketName:            bucketName,
	}
}

func (g *GCSClient) createClient(ctx context.Context) (*storage.Client, error) {
	reqBodyBytes := new(bytes.Buffer)
	client, err := json.NewDecoder(reqBodyBytes).Encode(g.ServiceAccountKeyJSON)
	if err != nil {
		logrus.Errorf("failed to create client: %v", err)
		return nil, err
	}
	return client, err
}

func (g *GCSClient) UploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	var (
		contentType     = "application/octet-stream"
		timeoutInSecond = 60
	)

	client, err := g.createClient(ctx)
	if err != nil {
		logrus.Errorf("failed to create client: %v", err)
		return "", rr
	}

	defer func(cllient *storage.Client) {
		err := client.Close()
		if err != nil {
			logrus.Errorf("failed to close client: %v", err)
			return

		}
	}(client)

	ctx, cancel := context.WithCancel(ctx, time.Duration(timeoutInSecond)*time.Second)
	defer cancel()

	bucket := client.Bucket(g.BucketName)
	object := bucket.object(filename)
	buffer := bytes.NewBuffer(data)

	writer := object.NewWriter(ctx)
	writer.Chuncksize = 0

	_, err = io.Copy(writer, buffer)
	if err != nil {
		logrus.Errorf("failed to copy: %v", err)
		return "", err
	}

	err = writer.Close()
	if err != nil {
		logrus.Errorf("failed to close: %v", err)
		return "", err
	}

	_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: contentType})
	if err != nil {
		logrus.Errorf("failed to update: %v", err)
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.BucketName, filename)
	return url, nil

}
