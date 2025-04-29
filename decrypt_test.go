package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// dummyDecryptor is a simple implementation of the Decryptor interface.
type dummyDecryptor struct{}

func (d dummyDecryptor) Init() error {
	fmt.Println("Init dummyDecryptor")
	return nil
}

// Decrypt returns "decryptedValue" for any secret name unless
// the secret name is "error", in which case it returns an error.
func (d dummyDecryptor) Decrypt(secretName string) (string, error) {
	if secretName == "error" {
		return "", errors.New("decryption failed")
	}
	return "decryptedValue", nil
}

func TestDecryptSecrets_SimpleString(t *testing.T) {
	// Setup test data with a simple string containing a secret.
	data := Data{
		"message": "Hello $$secret$$, welcome!",
	}
	expected := Data{
		"message": "Hello decryptedValue, welcome!",
	}

	// Call the decryption function.
	err := decryptSecrets(&data, dummyDecryptor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Compare the processed data.
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestDecryptSecrets_NestedMap(t *testing.T) {
	// Setup test data with a nested map.
	data := Data{
		"user": Data{
			"password": "$$password$$",
			"email":    "user@example.com",
		},
	}
	expected := Data{
		"user": Data{
			"password": "decryptedValue",
			"email":    "user@example.com",
		},
	}

	err := decryptSecrets(&data, dummyDecryptor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestDecryptSecrets_Slice(t *testing.T) {
	// Setup test data with a slice that contains strings with secret patterns,
	// as well as a nested map and another slice.
	data := Data{
		"items": []any{
			Data{
				"key": "$$secret$$",
			},
			[]any{
				"start-$$secret$$-end",
			},
		},
	}
	expected := Data{
		"items": []any{
			Data{
				"key": "decryptedValue",
			},
			[]any{
				"start-decryptedValue-end",
			},
		},
	}

	err := decryptSecrets(&data, dummyDecryptor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestDecryptSecrets_Error(t *testing.T) {
	// Setup test data where decryption is expected to fail.
	data := Data{
		"secret": "$$error$$",
	}

	err := decryptSecrets(&data, dummyDecryptor{})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	// Optionally, you can check the error message if needed.
	expectedErr := "decryption failed"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}
