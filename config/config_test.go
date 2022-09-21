package config

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
	assert.EqualError(err, "Invalid Environment: prod, production, dev, development, develop, testing, test, Given: wrong")
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
	assert.EqualError(err, "Invalid Configuration Type: JSON, YAML, TOML or \"\"(empty string), Given: txt")
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
