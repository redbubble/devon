package domain

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestNewApp(t *testing.T) {
	// Looks for a configuration file in the source code base directory
	//
	viper.Set("source-code-base-dir", "./test-fixtures/src")

	expectedSourceDir := filepath.Join("test-fixtures", "src", "foo")

	a, err := NewApp("foo", "development")

	if err != nil {
		t.Errorf("NewApp returned an unexpected error: %v\n", err)
	}

	if a.SourceDir != expectedSourceDir {
		t.Errorf("Expected '%s', got '%s'", expectedSourceDir, a.SourceDir)
	}

	// Reads the correct set of modes from the config file
	expectedModeNames := []string{"development", "unguessable"}

	a, err = NewApp("foo", "unguessable")

	if err != nil {
		t.Errorf("NewApp returned an unexpected error: %v\n", err)
	}

	if len(a.Config.Modes) != len(expectedModeNames) {
		t.Errorf("Expected %d modes, got %d.", len(expectedModeNames), len(a.Config.Modes))
	}

	for _, expectedName := range expectedModeNames {
		var found bool
		for name, _ := range a.Config.Modes {
			if name == expectedName {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Did not find expected mode '%s' when reading config", expectedName)
		}
	}

	// Each Mode has the correct working directory (defaulting to the App's SourceDir if unset)
	expectedWorkingDirs := map[string]string{
		"development": filepath.Join("test-fixtures", "src", "foo"),
		"unguessable": filepath.Join("test-fixtures", "src", "foo", "fergus"),
	}

	a, err = NewApp("foo", "development")

	if err != nil {
		t.Errorf("NewApp returned an unexpected error: %v\n", err)
	}

	if len(a.Config.Modes) != len(expectedWorkingDirs) {
		t.Errorf("Expected %d modes, got %d.", len(expectedWorkingDirs), len(a.Config.Modes))
	}

	for modeName, expectedDir := range expectedWorkingDirs {
		actualDir := a.Config.Modes[modeName].WorkingDir

		if expectedDir != actualDir {
			t.Errorf("Expected working dir for '%s' mode to be '%s', but it was '%s'.", modeName, expectedDir, actualDir)
		}
	}
}
