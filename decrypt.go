package main

import (
	"regexp"
	"strings"
)

// Decryptor interface for decrypting secrets
type Decryptor interface {
	Decrypt(secretName string) (string, error)
}

// decryptSecrets recursively decrypts values in a Data map.
func decryptSecrets(data *Data, decryptor Decryptor) error {
	secretPattern := regexp.MustCompile(`\$\$(.*?)\$\$`)
	for k, v := range *data {
		switch value := v.(type) {
		case string:
			// Handle string values.
			matches := secretPattern.FindAllStringSubmatch(value, -1)
			newStr := value
			for _, match := range matches {
				secretName := match[1]
				decryptedValue, err := decryptor.Decrypt(secretName)
				if err != nil {
					return err
				}
				newStr = strings.ReplaceAll(newStr, match[0], decryptedValue)
			}
			(*data)[k] = newStr
		case map[string]any:
			// Recursively process nested maps.
			nestedData := Data(value)
			if err := decryptSecrets(&nestedData, decryptor); err != nil {
				return err
			}
			(*data)[k] = nestedData
    // used for tests, basically the same as map[string]any
    case Data:
			// Process nested maps of the custom Data type.
			if err := decryptSecrets(&value, decryptor); err != nil {
				return err
			}
			(*data)[k] = value
		case []any:
			// Process slices using a helper function.
			newList, err := decryptSecretsInSlice(value, decryptor)
			if err != nil {
				return err
			}
			(*data)[k] = newList
		}
	}
	return nil
}

// decryptSecretsInSlice recursively processes entries in a slice.
func decryptSecretsInSlice(list []any, decryptor Decryptor) ([]any, error) {
	secretPattern := regexp.MustCompile(`\$\$(.*?)\$\$`)
	for i, v := range list {
		switch value := v.(type) {
		case string:
			matches := secretPattern.FindAllStringSubmatch(value, -1)
			newStr := value
			for _, match := range matches {
				secretName := match[1]
				decryptedValue, err := decryptor.Decrypt(secretName)
				if err != nil {
					return nil, err
				}
				newStr = strings.ReplaceAll(newStr, match[0], decryptedValue)
			}
			list[i] = newStr
		case map[string]any:
			nestedData := Data(value)
			if err := decryptSecrets(&nestedData, decryptor); err != nil {
				return nil, err
			}
			list[i] = nestedData
		case Data:
			nestedData := Data(value)
			if err := decryptSecrets(&nestedData, decryptor); err != nil {
				return nil, err
			}
			list[i] = nestedData
		case []any:
			// Call recursively for nested slices.
			newList, err := decryptSecretsInSlice(value, decryptor)
			if err != nil {
				return nil, err
			}
			list[i] = newList
		}
	}
	return list, nil
}
