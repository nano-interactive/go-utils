package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadLinesSuccess(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fPath := filepath.Join(dir, "file.json")
	assert := require.New(t)
	data := []string{"Test 1", "Test 2"}
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		t.Errorf("Failed to open file %s: %v", fPath, err)
		t.FailNow()
	}

	for _, item := range data {
		_, _ = file.Write([]byte(item))
		_, _ = file.Write([]byte("\n"))
	}

	lines := ReadLines(t, file)
	assert.Equal(data, lines)
}

func TestReadJsonLineSuccess(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	dir := t.TempDir()
	fPath := filepath.Join(dir, "file.json")

	type Product struct {
		Name  string
		Price float64
	}

	product := Product{
		Name:  "Product 1",
		Price: 555.333,
	}
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0o777)

	assert.NoError(err)

	bytes, _ := json.Marshal(product)
	_, _ = file.Write(bytes)
	_, _ = file.Write([]byte("\n"))

	productsFile, done := ReadJsonLine[Product](t, file)()

	assert.True(done)
	assert.Equal(productsFile, product)
}

func TestReadJsonDataSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	dir := t.TempDir()

	fName := filepath.Join(dir, "data.json")

	type Product struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	products := []Product{
		{Name: "Product 1", Price: 555.333},
		{Name: "Product 2", Price: 555.333},
	}
	file, err := os.OpenFile(fName, os.O_CREATE|os.O_RDWR, 0o777)
	assert.NoError(err)

	for _, product := range products {
		bytes, _ := json.Marshal(product)
		_, _ = file.Write(bytes)
		_, _ = file.Write([]byte("\n"))
	}

	productsFile := ReadJsonData[Product](t, file)

	assert.NotEmpty(productsFile)
	assert.Len(productsFile, 2)
	assert.Equal(productsFile, products)
}

func TestTempJsonLogFile(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	file := TempJsonLogFile(t, WithName("log.json"), Create(), ReadWrite())

	assert.Contains(file.Name(), "log.json")
}
