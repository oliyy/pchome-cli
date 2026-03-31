package config

import (
	"os"
	"path/filepath"
)

func writeFixtureConfig(path, fixtureName string) error {
	data, err := os.ReadFile(filepath.Join("testdata", fixtureName))
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
