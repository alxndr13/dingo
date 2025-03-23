package decrypt

import (
	"testing"
)

func TestExampleDecryptor(t *testing.T) {
	decryptor := &ExampleDecryptor{}

	secretName := "anySecretName"
	decrypted, err := decryptor.Decrypt(secretName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "decryptedValue"
	if decrypted != expected {
		t.Errorf("expected %q, got %q", expected, decrypted)
	}
}
