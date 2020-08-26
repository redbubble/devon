package domain

import (
	"path/filepath"
	"testing"
)

func TestNewApp(t *testing.T) {
	// Looks for a configuration file in the source code base directory
	//
	SourceCodeBaseDir = "./test-fixtures/src"

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

}
