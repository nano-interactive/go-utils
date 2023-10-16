package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNoArgs(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	cfg, err := New(Config{})

	// Assert
	assert.Error(err)
	assert.Nil(cfg)
}

func TestNewWrongEnvironment(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	cfg, err := New(Config{Env: "wrong"})

	// Assert
	assert.Error(err)
	assert.Nil(cfg)
	assert.EqualError(err, "invalid Environment: prod, production, dev, development, develop, testing, test")
}

func TestNewWrongParseType(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	cfg, err := New(Config{
		Type: "txt",
	})

	// Assert
	assert.Error(err)
	assert.Nil(cfg)
	assert.EqualError(err, "invalid Configuration Type: JSON, YAML, TOML or \"\"(empty string)")
}

func TestNewSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	file, _ := os.Create("test.json")
	type testConfig struct {
		Key string `json:"key"`
	}
	config := testConfig{Key: "test"}
	val, _ := json.Marshal(config)
	_, _ = file.Write(val)
	t.Cleanup(func() {
		_ = os.Remove(file.Name())
	})
	// Act
	cfg, err := New(Config{
		Env:   "testing",
		Name:  "test",
		Type:  "json",
		Paths: []string{"."},
	})

	// Assert
	assert.NoError(err)
	assert.NotNil(cfg)
	assert.Equal("test", cfg.Get("key"))
}

func TestProduction(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	file, _ := os.Create("test-production.json")
	type testConfig struct {
		Key string `json:"key"`
	}
	config := testConfig{Key: "test"}
	val, _ := json.Marshal(config)
	_, _ = file.Write(val)
	t.Cleanup(func() {
		_ = os.Remove(file.Name())
	})
	// Act
	cfg, err := New(Config{
		Env:  "production",
		Name: "test-production",
		Type: "json",
	})

	// Assert
	assert.NoError(err)
	assert.NotNil(cfg)
	assert.Equal("test", cfg.Get("key"))
}
