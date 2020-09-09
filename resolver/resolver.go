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
	for _, depName := range depChain {
		if depName == app.Name {
			depChainStr := strings.Join(append(depChain, app.Name), "->")
			err := fmt.Errorf("Circular dependency when resolving dependencies: %s", depChainStr)
			return []domain.App{}, err
		}
	}

	// Does the list of apps already include this app?
	// If so, we need to compare the modes.
	//   If the modes are the same, then this dependency is already satisfied.
	//    There's nothing to do, so we can return now.
	//   If the modes are different, we have conflicting dependencies.
	//    Return an error.
	for j := 0; j < len(apps); j++ {
		if apps[j].Name == app.Name {
			if apps[j].Mode.Name == app.Mode.Name {
				// The app we want, in the mode we want, is
				// already in the list. Therefore, there's
				// nothing to do.
				return apps, nil
			} else {
				err := fmt.Errorf("Dependency conflict when resolving dependencies: %s in '%s' mode conflicts with %s in '%s' mode.",
					app.Name,
					app.Mode.Name,
					apps[j].Name,
					apps[j].Mode.Name)

				return []domain.App{}, err
			}
		}
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
