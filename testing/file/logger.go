package file

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type fileConfig struct {
	name        string
	dir         string
	permissions os.FileMode
	flags       int
}

type Modifier func(*fileConfig)

func WithPermissions(perms os.FileMode) Modifier {
	return func(c *fileConfig) {
		c.permissions = perms
	}
}

func WithName(name string) Modifier {
	return func(c *fileConfig) {
		c.name = name
	}
}

func WithDirectory(dir string) Modifier {
	return func(fc *fileConfig) {
		fc.dir = dir
	}
}

func Append() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_APPEND
	}
}

func ReadOnly() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_RDONLY
	}
}

func ReadWrite() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_RDWR
	}
}

func WriteOnly() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_WRONLY
	}
}

func Truncate() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_TRUNC
	}
}

func Create() Modifier {
	return func(c *fileConfig) {
		c.flags |= os.O_CREATE
	}
}

// io.Reader io.Write io.Closer io.Seeker fmt.Stringer error
// net.Addr net.Conn net.Listener net.Dialer

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

func ReadLines(t testing.TB, input io.Reader) []string {
	t.Helper()

	bytes := ReadLinesInBytes(t, input)

	lines := make([]string, 0, len(bytes))

	for _, b := range bytes {
		lines = append(lines, string(b))
	}

	return lines
}

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
