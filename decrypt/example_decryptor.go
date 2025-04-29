package decrypt

type ExampleDecryptor struct{}

func NewExampleDecryptor() *ExampleDecryptor {
	return &ExampleDecryptor{}
}

func (d *ExampleDecryptor) Init() error {
	return nil
}
func (d *ExampleDecryptor) Decrypt(secretName string) (string, error) {
	// Here, we'll just return a dummy value for demonstration
	return "decryptedValue", nil
}
