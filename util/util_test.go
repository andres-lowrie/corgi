package util

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var tmpDir string

func checkErrors(t *testing.T, caseNum int, expectedError, gotError error) {
	if expectedError != nil && gotError == nil {
		t.Errorf("Case %d: Should have returned an error: %v", caseNum, gotError)
	}

	if expectedError == nil && gotError != nil {
		t.Errorf("Case %d: Should not have returned an error: %v", caseNum, gotError)
	}
}

func TestMain(m *testing.M) {
	tmpDir = os.Getenv("TMPDIR") + "/corgi_test"
	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func TestLoadJsonDataFromFile(t *testing.T) {
	tests := []struct {
		filePath string
		object   interface{}
		err      error
	}{
		// Happy Path
		{"testdata/a_json_file.json", new(struct{ Foo string }), nil},
		// It should bubble up errors
		{"testdata/no_such_file", new(interface{}), fmt.Errorf("")},
		{"testdata/bad_json_file.json", new(interface{}), fmt.Errorf("")},
	}

	for i, tt := range tests {
		got := LoadJsonDataFromFile(tt.filePath, &tt.object)
		checkErrors(t, i, tt.err, got)
	}
}

func TestScan(t *testing.T) {
	// Setup
	mockHistoryFileLoc := "testdata/history_file"

	tests := []struct {
		prompt      string
		defaultInp  string
		historyFile string
		expected    string
		err         error
	}{
		// Happy path
		{"", "def", mockHistoryFileLoc, "def", nil},
		// It should handle bad inputs
		{"", "", mockHistoryFileLoc, "", fmt.Errorf("")},
		{"", " ", mockHistoryFileLoc, "", fmt.Errorf("")},
		{"\a", "", mockHistoryFileLoc, "", fmt.Errorf("")},
		// It should handle multiline commands
		{"", "\\", mockHistoryFileLoc, "", fmt.Errorf("")},
	}

	for i, tt := range tests {
		got, err := Scan(tt.prompt, tt.defaultInp, tt.historyFile)

		if got != tt.expected {
			t.Errorf("Case %d: Exptected `%+v`, got `%+v`", i, tt.expected, got)
		}
		checkErrors(t, i, tt.err, err)
	}

	// Teardown
	os.Remove(mockHistoryFileLoc)
}

func TestExecute(t *testing.T) {
	tests := []struct {
		setup   func()
		command string
		err     error
	}{
		// Happy path
		{func() {}, "echo blah", nil},
		// It should default to "sh" if no shell available
		{func() { os.Unsetenv("SHELL") }, "echo blah", nil},
		// It should bubble up errors for bad commands
		{func() { os.Setenv("SHELL", "foo") }, "echo blah", fmt.Errorf("")},
		{func() {}, "nosuchcommand", fmt.Errorf("")},
	}

	// Need to pass something that implements Reader and Writer but we don't need
	// to actually read from it so we'll just create a buffer and use that for
	// all the tests
	for i, tt := range tests {
		tt.setup()

		var buf bytes.Buffer
		err := Execute(tt.command, &buf, &buf)
		checkErrors(t, i, tt.err, err)
	}
}

func TestGetOrCreatePath(t *testing.T) {
	tests := []struct {
		loc   string
		perm  os.FileMode
		isDir bool
		err   error
	}{
		// Happy path
		{"wont/exist", 0755, false, nil},
		{"newfile", 0644, false, nil},
		{"dir", 0644, true, nil},
		// Should bubble up errors
		{"some/other/path", 0644, false, fmt.Errorf("")},
	}

	for i, tt := range tests {
		path := fmt.Sprintf("%s/%s", tmpDir, tt.loc)
		err := GetOrCreatePath(path, tt.perm, tt.isDir)
		spew.Dump(err)
		checkErrors(t, i, tt.err, err)
	}
}
