package file

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// type fileConfig
// Structure contains basic file options: name of file, directory to be written to, permissions and flags
type fileConfig struct {
	name        string
	dir         string
	permissions os.FileMode
	flags       int
}

// type Modifier
// Delegates functions for a given file
type Modifier func(*fileConfig)

// Sets file premissions for testing purposes
func WithPermissions(perms os.FileMode) Modifier {
	return func(c *fileConfig) {
		c.permissions = perms
	}
}

// Sets file name for testing purposes
func WithName(name string) Modifier {
	return func(c *fileConfig) {
		c.name = name
	}
}

// Sets directory for testing purposes
func WithDirectory(dir string) Modifier {
	return func(fc *fileConfig) {
		fc.dir = dir
	}
}

// Sets Append flag to a file for testing purposes
func Append() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_APPEND
	}
}

// Marks file as read only for testing purposes
func ReadOnly() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_RDONLY
	}
}

// Marks file as read-write for testing purposes
func ReadWrite() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_RDWR
	}
}

// Marks file as write-only for testing purposes
func WriteOnly() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_WRONLY
	}
}

// Marks file as truncate for testing purposes
func Truncate() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_TRUNC
	}
}

// Sets Create flag to a file for testing purposes
func Create() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_CREATE
	}
}

// Returns T and bool depending on failing or not for a given file for testing purposes
func ReadJsonLine[T any](t testing.TB, input io.Reader) func() (T, bool) {
	t.Helper()

	if seek, ok := input.(io.Seeker); ok {
		if _, err := seek.Seek(0, io.SeekStart); err != nil {
			t.Log(err)
			t.FailNow()
		}
	}

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	return func() (T, bool) {
		t.Helper()
		var decode T

		if !scanner.Scan() {
			return decode, false
		}

		if err := json.Unmarshal(scanner.Bytes(), &decode); err != nil {
			t.Log(err)
			t.FailNow()
		}

		return decode, true
	}
}

// Returns []T for given file for testing purposes
func ReadJsonData[T any](t testing.TB, input io.Reader) []T {
	t.Helper()

	var decode T
	storage := make([]T, 0, 100)
	lines := ReadLinesInBytes(t, input)

	for _, line := range lines {
		if err := json.Unmarshal(line, &decode); err != nil {
			t.Log(err)
			t.FailNow()
		}

		storage = append(storage, decode)
	}

	return storage
}

// Returns log data in Byte format for single line for testing purposes
func ReadLinesInBytes(t testing.TB, input io.Reader) [][]byte {
	t.Helper()

	if seek, ok := input.(io.Seeker); ok {
		if _, err := seek.Seek(0, io.SeekStart); err != nil {
			t.Log(err)
			t.FailNow()
		}
	}

	lines := make([][]byte, 0, 100)

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	return lines
}

// Returns log data in String slice format from Byte format for testing purposes
func ReadLines(t testing.TB, input io.Reader) []string {
	t.Helper()

	bytes := ReadLinesInBytes(t, input)

	lines := make([]string, 0, len(bytes))

	for _, b := range bytes {
		lines = append(lines, string(b))
	}

	return lines
}

// Creates temporary JSON log file for testing purposes
func TempJsonLogFile(t testing.TB, modifiers ...Modifier) *os.File {
	t.Helper()

	cfg := fileConfig{
		name:  "test.json",
		flags: os.O_CREATE | os.O_RDWR,
		dir:   t.TempDir(),
	}

	for _, modifier := range modifiers {
		modifier(&cfg)
	}

	filePath := filepath.Join(cfg.dir, cfg.name)

	// #nosec 304
	file, err := os.OpenFile(filePath, cfg.flags, cfg.permissions)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Log(err)
			t.FailNow()
		}

		if err := os.Remove(filePath); err != nil {
			t.Log(err)
			t.FailNow()
		}
	})

	return file
}
