package config

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/gofor-little/env"
)

func Load(filepath string) error {
	rootDir := findProjectRoot()
	path := path.Join(rootDir, filepath)
	return envLoader(path)
}

func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		// Check if go.mod exists in the current directory traversal
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}

		// Move up one directory
		parent := filepath.Dir(wd)
		if parent == wd {
			// Reached the system root directory without finding go.mod
			break
		}
		wd = parent
	}

	return "."
}

func envLoader(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	return env.Load(filepath)
}

func GetEnv[T any](key string, defaultValue T) T {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var target T

	// Handle time.Duration explicitly by its TypeOf
	tType := reflect.TypeOf(target)

	if tType == reflect.TypeOf(time.Duration(0)) {
		if d, err := time.ParseDuration(valStr); err == nil {
			reflect.ValueOf(&target).Elem().Set(reflect.ValueOf(d))
			return target
		}
		return defaultValue
	}

	kind := reflect.TypeOf(target).Kind()

	switch kind {
	case reflect.String:
		reflect.ValueOf(&target).Elem().SetString(valStr)
		return target

	case reflect.Int:
		if i, err := strconv.Atoi(valStr); err == nil {
			reflect.ValueOf(&target).Elem().SetInt(int64(i))
			return target
		}

	case reflect.Bool:
		if b, err := strconv.ParseBool(valStr); err == nil {
			reflect.ValueOf(&target).Elem().SetBool(b)
			return target
		}
	}

	return defaultValue
}
