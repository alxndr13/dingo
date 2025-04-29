package decrypt

import (
	"context"
	"errors"
	"testing"

	"github.com/googleapis/gax-go/v2"

	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Mocked Google Secret Manager Client
type mockSecretManagerClient struct {
	accessSecretVersionFunc func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error)
	closeFunc               func() error
}

func (m *mockSecretManagerClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return m.accessSecretVersionFunc(ctx, req)
}

func (m *mockSecretManagerClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestGoogleDecryptor_Decrypt(t *testing.T) {
	const secretName = "projects/my-project/secrets/my-secret/versions/latest"
	const secretValue = "supersecretvalue"

	mockClient := &mockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			if req.Name != secretName {
				t.Errorf("unexpected secret name: got %s, want %s", req.Name, secretName)
			}
			return &secretmanagerpb.AccessSecretVersionResponse{
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte(secretValue),
				},
			}, nil
		},
	}

	decryptor := NewGoogleDecryptor()

	decryptor.client = mockClient

	got, err := decryptor.Decrypt(secretName)
	if err != nil {
		t.Fatalf("Decrypt() error = %v, want nil", err)
	}
	if got != secretValue {
		t.Errorf("Decrypt() = %q, want %q", got, secretValue)
	}
}

func TestGoogleDecryptor_Decrypt_Error(t *testing.T) {
	mockClient := &mockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			return nil, errors.New("some error")
		},
	}

	decryptor := NewGoogleDecryptor()
	decryptor.client = mockClient

	_, err := decryptor.Decrypt("projects/my-project/secrets/my-secret/versions/latest")
	if err == nil {
		t.Fatal("Decrypt() error = nil, want error")
	}
}
