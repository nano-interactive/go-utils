package file

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadLinesSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	const fileName = "file.json"
	assert := require.New(t)
	data := []string{"Test 1", "Test 2"}
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0o777)
	for _, item := range data {
		_, _ = file.Write([]byte(item))
		_, _ = file.Write([]byte("\n"))
	}

	// Act
	lines := ReadLines(t, file)

	// Assert
	assert.Equal(data, lines)

	// Cleanup
	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Log(err)
			t.FailNow()
		}

		if err := os.Remove(fileName); err != nil {
			t.Log(err)
			t.FailNow()
		}
	})
}

func TestReadJsonLineSuccess(t *testing.T) {
	// Arrange
	const fileName = "products.json"
	t.Parallel()
	assert := require.New(t)
	type Product struct {
		Name  string
		Price float64
	}

	product := Product{
		Name:  "Product 1",
		Price: 555.333,
	}
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0o777)
	bytes, _ := json.Marshal(product)
	_, _ = file.Write(bytes)
	_, _ =file.Write([]byte("\n"))

	// Act
	productsFile, _ := ReadJsonLine[Product](t, file)()
	// Assert
	assert.Equal(productsFile, product)

	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Log(err)
			t.FailNow()
		}

		if err := os.Remove(fileName); err != nil {
			t.Log(err)
			t.FailNow()
		}
	})
}

func TestReadJsonDataSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	const fileName = "data.json"

	type Product struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	products := []Product{
		{Name: "Product 1", Price: 555.333},
		{Name: "Product 2", Price: 555.333},
	}
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0o777)

	for _, product := range products {
		bytes, _ := json.Marshal(product)
		file.Write(bytes)
		file.Write([]byte("\n"))
	}

	// Act
	productsFile := ReadJsonData[Product](t, file)

	// Assert
	assert.NotEmpty(productsFile)
	assert.Len(productsFile, 2)
	assert.Equal(productsFile, products)

	// Cleanup
	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Log(err)
			t.FailNow()
		}

		if err := os.Remove(fileName); err != nil {
			t.Log(err)
			t.FailNow()
		}
	})
}

func TestTempJsonLogFile(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)

	// Act
	file := TempJsonLogFile(t, WithName("log.json"), Create(), ReadWrite())

	// Assert
	assert.Contains(file.Name(), "log.json")
}
