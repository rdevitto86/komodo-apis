package config

import (
	"os"
	"sync"
)

type Config struct {
	mu   sync.RWMutex
	data map[string]string
}

var (
	once     sync.Once
	instance *Config
)

func init() {
	once.Do(func() {
		instance = &Config{ data: make(map[string]string) }
	})
}

// GetConfigValue checks local in-memory config first, then falls back to environment variable
func GetConfigValue(key string) string {
	if key == "" || instance == nil { return "" }
	instance.mu.RLock()
	val := instance.data[key]
	instance.mu.RUnlock()
	if val != "" { return val }
	return os.Getenv(key)
}

// SetConfigValue sets value in local in-memory config
func SetConfigValue(key, value string) {
	if value == "" || key == "" || instance == nil { return }
	instance.mu.Lock()
	instance.data[key] = value
	instance.mu.Unlock()
}

// DeleteLocalConfig removes value from local in-memory config only
func DeleteConfigValue(key string) {
	if key == "" || instance == nil { return }
	instance.mu.Lock()
	delete(instance.data, key)
	instance.mu.Unlock()
}
