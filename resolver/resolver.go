package resolver

import (
	"fmt"

	"github.com/redbubble/devon/domain"
)

func Add(apps []domain.App, app domain.App) ([]domain.App, error) {
	starting_apps, err := add(apps, app, make([]string, 0, len(apps)))

	if err != nil {
		return nil, err
	}

	return starting_apps, nil
}

func add(apps []domain.App, app domain.App, dep_chain []string) ([]domain.App, error) {
	// Does the dependency chain already include this app?
	// If so, we've got a circular dependency, so return an error.


	// Does the list of apps already include this app?
	// If so, we need to compare the modes.
	//   If the modes are the same, then this dependency is already satisfied.
	//    There's nothing to do, so we can return now.
	//   If the modes are different, we have conflicting dependencies.
	//    Return an error.


	// OK, we've gotten this far without errors.
	// Add the app to the list


	// Recursively add this app's dependencies

}
