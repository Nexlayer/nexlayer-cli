package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// CreateTempFile creates a temporary file with the given content
func CreateTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "nexlayer-test-*.yaml")
	require.NoError(t, err)

	_, err = tmpfile.Write([]byte(content))
	require.NoError(t, err)

	err = tmpfile.Close()
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

// ReadTestData reads test data from the testdata directory
func ReadTestData(t *testing.T, filename string) []byte {
	t.Helper()
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

// CompareJSON compares two JSON strings for equality
func CompareJSON(t *testing.T, expected, actual string) {
	t.Helper()
	var expectedJSON, actualJSON interface{}

	err := json.Unmarshal([]byte(expected), &expectedJSON)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(actual), &actualJSON)
	require.NoError(t, err)

	require.Equal(t, expectedJSON, actualJSON)
}

// CaptureOutput captures stdout during test execution
func CaptureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()

	w.Close()
	os.Stdout = old
	return <-outC
}
