# Devon, for starting systems in dev

_Because it's time to do dev on our stuff!_

## Usage

### In an application repo

```
devon start
```

starts the local application (whose repo is the current working directory) in Development mode, with its dependencies in Dependency mode (see [Modes](#modes)).

### From another directory

```
devon start <repo-name>
```

starts the named application in Development mode, with its dependencies in Dependency mode (see [Modes](#modes)).

### Options

`--mode`, `-m`
: Specify the mode. Default: `development`

`--verbose`, `-v`
: Print some more information about what's happening.

## Modes

You can run applications in various modes, to suit particular dev/test scenarios. Want to run the web app without the background workers? Make a mode for it. Want to run both? Another mode. Using placeholder data instead of calling out to a dependency? You guessed it: add a mode!

See some examples in `example.yaml` in the root of this repo.

### Development mode

Runs the application so that you can make changes to the code and see them in your running application. Usually this means running the application natively on your laptop, although there are exceptions. The exact implementation depends on that repo's `dev/dev.sh`.

### Dependency mode

Runs the application, but doesn't pay attention to code changes. Usually, this works by running a Docker container, though the exact implementation depends on that application's `dev/up.sh`.

### Custom modes

Each application can specify custom modes, to enable or disable specific functions. For example, an application with both a synchronous frontend and an async backend worker might have the worker disabled by default, and use a custom mode to enable it, or to run it in isolation.

## Important assumptions

### All app repos live in ~/src

It'd be good to allow for overriding this, but it's not done yet.

### One app should never be started in 2 modes at the same time

I (Lucas) can't think of a sane use case for that, so I'm pretty sure this assumption is safe.

### Applications run in the background

If Devon starts applications by running commands naively (e.g. without spawning a child process), then those applications can interrupt the starting process by running in the foreground. Because they never exit, subsequent applications wouldn't be able to start.

### All dependency relationships are the same

At the moment, a dependency relationship means "A depends on B to work properly, but A can start even if B isn't running." This is optimistic, since there is nothing stopping A from crashing if B is not available at startup.

This probably needs fixing with an actual dependency graph implementation -- the current implementation just builds a list that includes all the dependencies, with no information about which apps depend on which others. Hence, there's no reliable way to figure out what to start first.

## Development

**[Trello board](https://trello.com/b/MsxE9Nw6/devon-the-dev-application-starter)**
