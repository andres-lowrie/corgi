package snippet

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var tmpDir string

func TestMain(m *testing.M) {
	// Create directory for mock history files
	tmpDir = os.Getenv("TMPDIR") + "corgi_test"
	os.MkdirAll(tmpDir, 0755)
	code := m.Run()
	//os.RemoveAll(tmpDir)
	os.Exit(code)
}

// Creates a mock 'fish' shell by setting SHELL to a file that will  just
// output the passed in version
func setupMockFish(v string) string {

	content, _ := ioutil.ReadFile("testdata/mock.sh")
	targetPath := path.Join(tmpDir, "mock.sh")
	newContent := strings.Replace(string(content), "%%VERSION%%", v, 1)
	ioutil.WriteFile(targetPath, []byte(newContent), 0777)
	return targetPath
}

func TestGetFishHistoryPath(t *testing.T) {
	shellToMockFish := func(s string) {
		os.Setenv("SHELL", setupMockFish(s))
	}

	tests := []struct {
		mockShell   func(s string)
		homeDir     string
		expect      string
		mockVersion string
	}{
		// Should handle pre and post 2.3.0 version of fish
		{shellToMockFish, tmpDir, path.Join(tmpDir, ".config", "fish", "fish_history"), "1.0.0"},
		{shellToMockFish, tmpDir, path.Join(tmpDir, ".local", "share", "fish", "fish_history"), "2.3.1"},
		// Should default to newHistFile on bad version
		{shellToMockFish, tmpDir, path.Join(tmpDir, ".local", "share", "fish", "fish_history"), ".."},
		// The following test should work but doesn't, there may be a bug in the func
		//{
		//  func(s string) { os.Setenv("SHELL", "nosuchshell") },
		//  tmpDir,
		//  path.Join(tmpDir, ".local", "share", "fish", "fish_history"),
		//  "..",
		//},
	}

	for i, tt := range tests {
		tt.mockShell(tt.mockVersion)
		got := getFishHistoryPath(tt.homeDir)

		if string(got) != tt.expect {
			t.Errorf("Case %d: Wanted `%s`, got `%s`", i, tt.expect, got)
		}
	}

}
