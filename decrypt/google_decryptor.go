package decrypt

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
)

type SecretManagerClient interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

type GoogleDecryptor struct {
	client SecretManagerClient
}

func NewGoogleDecryptor() *GoogleDecryptor {
	return &GoogleDecryptor{}
}

func (d *GoogleDecryptor) Init() error {
	ctx := context.Background()

	// Create the Secret Manager client using ADC
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	d.client = client
	return nil
}

func (d *GoogleDecryptor) Decrypt(secretName string) (string, error) {
	ctx := context.Background()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := d.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	secretData := string(result.Payload.Data)

	return secretData, nil
}
