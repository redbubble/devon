package resolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/redbubble/devon/domain"
)

func TestAdd(t *testing.T) {
	workdir, _ := os.Getwd()
	domain.DefaultSourceCodeBase = filepath.Join(workdir, "test-fixtures", "src")

	// Adds the given app to the list of apps to be started
	apps := []domain.App{}
	app := domain.App{
		Name: "testfenster",
		SourceDir: "",
		Config: domain.Config{},
		Mode: domain.Mode{
			Dependencies: make(map[string]string),
		},
	}

	actual, err := Add(apps, app)

	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

	expectedAppNames := []string{ app.Name }

	if len(actual) != len(expectedAppNames) {
		t.Errorf("Expected %d apps, got %d", len(expectedAppNames), len(actual))
	}

	for i, expectedAppName := range expectedAppNames {
		if actual[i].Name != expectedAppName {
			t.Errorf("Expected %s, got %v", expectedAppName, actual[i].Name)
		}
	}

	// Adds the given app's dependencies to the list of apps to be started
	dependencies := make(map[string]string)
	dependencies["first-dep"] = "dependency"

	appWithDependencies := domain.App{
		Name: "testfenster",
		SourceDir: "",
		Config: domain.Config{},
		Mode: domain.Mode{
			Dependencies: dependencies,
		},
	}

	actual, err = Add(apps, appWithDependencies)

	expectedAppNames = []string{
		"testfenster",
		"first-dep",
	}

	if len(actual) != len(expectedAppNames) {
		t.Errorf("Expected %d apps, got %d", len(expectedAppNames), len(actual))
	}

	for i, expectedAppName := range expectedAppNames {
		if actual[i].Name != expectedAppName {
			t.Errorf("Expected %s, got %v", expectedAppName, actual[i].Name)
		}
	}

	// Detects circular dependencies and returns an error
	dependencies = make(map[string]string)
	dependencies["kreis"] = "dependency"

	appWithCircularDependency := domain.App{
		Name: "kreis",
		SourceDir: "",
		Config: domain.Config{},
		Mode: domain.Mode{
			Dependencies: dependencies,
		},
	}

	_, expectedErr := Add(apps, appWithCircularDependency)

	if expectedErr == nil {
		t.Errorf("Expected a circular dependency to result in an error, but it didn't.")
	}

	// Detects conflicting dependencies
	apps = []domain.App{
		domain.App{
			Name: "testfenster",
			Mode: domain.Mode{
				Name: "dependency",
			},
		},
	}

	conflictingApp := domain.App{
		Name: "testfenster",
		SourceDir: "",
		Config: domain.Config{},
		Mode: domain.Mode{
			Name: "development",
		},
	}

	_, expectedErr = Add(apps, conflictingApp)

	if expectedErr == nil {
		t.Errorf("Expected a conflicting dependency (same name, different mode) to result in an error, but it didn't.")
	}
}
