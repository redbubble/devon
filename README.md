# Devon, for starting systems in dev

_Because it's time to do dev on our stuff!_

## Usage

### In an application repo

```
devon
```

starts the local application (whose repo is the current working directory) in Development mode, with its dependencies in Dependency mode (see [Modes](#modes)).

### From another directory

```
devon <repo-name>
```

starts the named application in Development mode, with its dependencies in Dependency mode (see [Modes](#modes)).

### Using a config file

```
devon -f <config-path>
```

### Options

`--mode`, `-m`
: Specify the mode. Default: `development`

`--file`, `-f`
: Use a configuration file that specifies a set of applications and modes to run.


## Modes

You can run applications in various modes, to suit particular dev/test scenarios. Want to run the web app without the background workers? Make a mode for it. Want to run both? Another mode. Using placeholder data instead of calling out to a dependency? You guessed it: add a mode!

See some examples in `example.yaml` in the root of this repo.

### Development mode

Runs the application so that you can make changes to the code and see them in your running application. Usually this means running the application natively on your laptop, although there are exceptions. The exact implementation depends on that repo's `dev/dev.sh`.

### Dependency mode

Runs the application, but doesn't pay attention to code changes. Usually, this works by running a Docker container, though the exact implementation depends on that application's `dev/up.sh`.

### Custom modes

Each application can specify custom modes, to enable or disable specific functions. For example, an application with both a synchronous frontend and an async backend worker might have the worker disabled by default, and use a custom mode to enable it, or to run it in isolation.

## Development

**[Trello board](https://trello.com/b/MsxE9Nw6/devon-the-dev-application-starter)**
