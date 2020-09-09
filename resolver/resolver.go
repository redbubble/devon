package resolver

import (
	"fmt"
	"strings"

	"github.com/redbubble/devon/domain"
)

func Add(apps []domain.App, app domain.App) ([]domain.App, error) {
	startingApps, err := add(apps, app, make([]string, 0, len(apps)))

	if err != nil {
		return nil, err
	}

	return startingApps, nil
}

func add(apps []domain.App, app domain.App, depChain []string) ([]domain.App, error) {
	var startingApps []domain.App

	// Does the dependency chain already include this app?
	// If so, we've got a circular dependency, so return an error.
	if checkCircularDependency(app.Name, depChain) {
		depChainStr := strings.Join(append(depChain, app.Name), "->")
		err := fmt.Errorf("Circular dependency when resolving dependencies: %s", depChainStr)

		return []domain.App{}, err
	}

	if conflict, conflictingApp := checkDependencyConflict(app, apps); conflict {
		err := fmt.Errorf("Dependency conflict when resolving dependencies: %s in '%s' mode conflicts with %s in '%s' mode.",
			app.Name,
			app.Mode.Name,
			conflictingApp.Name,
			conflictingApp.Mode.Name)

		return []domain.App{}, err
	}

	if checkAlreadyIncluded(app, apps) {
		return apps, nil
	}

	// All our validation checks have passed, so add the app to the list.
	startingApps = append(apps, app)

	// Recursively add this app's dependencies
	for _, dependency := range app.Mode.Dependencies {
		dep, err := domain.NewApp(dependency.AppName, dependency.ModeName)

		if err != nil {
			fmt.Printf("WARNING: %v\n", err)
			continue
		}

		startingApps, err = add(startingApps, dep, append(depChain, app.Name))

		if err != nil {
			return []domain.App{}, err
		}
	}

	return startingApps, nil
}

func checkCircularDependency(appName string, depChain []string) bool {
	for _, depName := range depChain {
		if depName == appName {
			return true
		}
	}

	return false
}

func checkDependencyConflict(app domain.App, apps []domain.App) (bool, domain.App) {
	for _, a := range apps {
		if a.Name == app.Name && a.Mode.Name != app.Mode.Name {
			return true, a
		}
	}

	return false, domain.App{}
}

func checkAlreadyIncluded(app domain.App, apps []domain.App) bool {
	for _, a := range apps {
		if a.Name == app.Name && a.Mode.Name == app.Mode.Name {
			return true
		}
	}

	return false
}
