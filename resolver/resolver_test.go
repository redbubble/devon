package resolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/redbubble/devon/domain"
)

func TestAdd(t *testing.T) {
	workdir, _ := os.Getwd()
	viper.Set("source-code-base-dir", filepath.Join(workdir, "test-fixtures", "src"))

	// Adds the given app to the list of apps to be started
	apps := []domain.App{}
	app := domain.App{
		Name:      "testfenster",
		SourceDir: "",
		Config:    domain.Config{},
		Mode: domain.Mode{
			Dependencies: make([]domain.Dependency, 0, 1),
		},
	}

	actual, err := Add(apps, app, []string{})

	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

	expectedAppNames := []string{app.Name}

	if len(actual) != len(expectedAppNames) {
		t.Errorf("Expected %d apps, got %d", len(expectedAppNames), len(actual))
	}

	for i, expectedAppName := range expectedAppNames {
		if actual[i].Name != expectedAppName {
			t.Errorf("Expected %s, got %v", expectedAppName, actual[i].Name)
		}
	}

	// Adds the given app's dependencies to the list of apps to be started
	dependencies := []domain.Dependency{
		domain.Dependency{
			AppName:  "first-dep",
			ModeName: "dependency",
		},
	}

	appWithDependencies := domain.App{
		Name:      "testfenster",
		SourceDir: "",
		Config:    domain.Config{},
		Mode: domain.Mode{
			Dependencies: dependencies,
		},
	}

	actual, err = Add(apps, appWithDependencies, []string{})

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
	dependencies = []domain.Dependency{
		domain.Dependency{
			AppName:  "kreis",
			ModeName: "dependency",
		},
	}

	appWithCircularDependency := domain.App{
		Name:      "kreis",
		SourceDir: "",
		Config:    domain.Config{},
		Mode: domain.Mode{
			Dependencies: dependencies,
		},
	}

	_, expectedErr := Add(apps, appWithCircularDependency, []string{})

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
		Name:      "testfenster",
		SourceDir: "",
		Config:    domain.Config{},
		Mode: domain.Mode{
			Name: "development",
		},
	}

	_, expectedErr = Add(apps, conflictingApp, []string{})

	if expectedErr == nil {
		t.Errorf("Expected a conflicting dependency (same name, different mode) to result in an error, but it didn't.")
	}

	// Declines to add apps from the skip list
	skippedApp := domain.App{
		Name: "skipme",
		Mode: domain.Mode{
			Name: "dependency",
		},
	}

	apps = []domain.App{}

	skip := []string{"skipme"}

	actualApps, err := Add(apps, skippedApp, skip)

	if err != nil {
		t.Error(err)
	}

	if len(actualApps) != 0 {
		t.Errorf("Expected an app appearing in the skip list to be skipped, but it wasn't.")
	}
}
