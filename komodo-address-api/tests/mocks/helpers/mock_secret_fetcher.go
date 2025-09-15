package helpers

import "fmt"

type MockSecretFetcher struct{}

func (m *MockSecretFetcher) GetSecret(name string) (string, error) {
    if name == "missing" {
        return "", fmt.Errorf("secret not found")
    }
    return "mocked-secret", nil
}
