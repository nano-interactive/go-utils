package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"unicode"
	"unsafe"
)

// #nosec G103
// Returns a byte pointer without allocation
func UnsafeBytes(s string) []byte {
	var bs []byte

	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len

	return bs
}

// #nosec G103
// Returns a string pointer without allocation
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Returns bool value of HTTP status code if its failure or success
func IsSuccess(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

// Provides a way of testing if type is integer
// Returns bool value depending if its integer or not
func IsInt(s string) bool {
	if len(s) > 0 && s[0] == '0' {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

// Returns ENV string
func Getenv(env string, def ...string) string {
	item, exists := os.LookupEnv(env)

	if !exists && len(def) > 0 {
		return def[0]
	}

	return item
}

// Returns random string of given length
func RandomString(n int32) string {
	buffer := make([]byte, n)

	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(buffer)
}

// Returns absolute path string for a given directory and error if directory doesent exist
func GetAbsolutePath(path string) (string, error) {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)

		if err != nil {
			return "", err
		}

		return path, nil
	}

	return path, err
}

// Provides a way of creating a directory from a path
// Returns created directory path and error if fails
func CreateDirectoryFromFile(path string, perm fs.FileMode) (string, error) {
	p, err := GetAbsolutePath(path)
	if err != nil {
		return "", err
	}

	directory := filepath.Dir(p)

	if err := os.MkdirAll(directory, perm); err != nil {
		return "", err
	}

	return p, nil
}

// Creates file for given directory with flags and permissions for directory and file
// Returns file instance and error if it fails
func CreateFile(path string, flags int, dirMode, mode fs.FileMode) (file *os.File, err error) {
	path, err = CreateDirectoryFromFile(path, dirMode|fs.ModeDir)

	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		//#nosec G304
		file, err = os.Create(path)

		if err != nil {
			return nil, err
		}

		if err = file.Chmod(mode); err != nil {
			return nil, err
		}

		if err = file.Close(); err != nil {
			return nil, err
		}
	}

	//#nosec G304
	file, err = os.OpenFile(path, flags, mode)

	return
}

// Creates write only appendable file with permission 0o744 for directory and file for given path
func CreateLogFile(path string) (file *os.File, err error) {
	file, err = CreateFile(path, os.O_WRONLY|os.O_APPEND, 0o744, fs.FileMode(0o744)|os.ModeAppend)

	return
}

// Provides a way of checking if file exists. Returns bool
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// Provides a way of creating directory with permissions for a given path
// Returns string path of created directory and error if fails
func CreateDirectory(path string, perm fs.FileMode) (string, error) {
	p, err := GetAbsolutePath(path)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(p, perm); err != nil {
		return "", err
	}

	return p, nil
}

// Provides a way of copying bytes into new byte slice
// Returns byte slice
func CopyBytes(input []byte) []byte {
	c := make([]byte, len(input))

	copy(c, input)

	return c
}

var pixel []byte = nil

const image = "R0lGODlhAQABAIAAAP///wAAACH5BAAAAAAALAAAAAABAAEAAAICRAEAOw=="

func init() {
	if pixel == nil {
		var err error
		pixel, err = base64.StdEncoding.DecodeString(image)

		if err != nil {
			panic(err)
		}
	}
}

func GetBrokenImageBytes() []byte {
	return pixel
}
