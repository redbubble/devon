# Devon, for starting systems in dev

_Because it's time to do dev on our stuff!_

## Installation

### Homebrew

If you're on a Mac, you can use Homebrew to install Devon:

```bash
brew tap redbubble/devon
brew install devon
```

### Download release from GitHub

You can download pre-compiled binaries from the [Releases](https://github.com/redbubble/devon/releases) page. There are OSX and Linux archives for download, as well as RPM and DEB packages if that's your flavour.

### Compile from source

You can of course clone this repo using `git clone git@github.com:redbubble/devon` and build Devon yourself if you wish.

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

In general, most apps should have at least two basic modes:

1. `development`, which is used when you are actively making changes to an application, and
2. `dependency`, which is used when you need an app to be running (e.g. to satisfy a runtime dependency of an app you're working on), but you aren't planning on making changes to it.

See some examples in `example.yaml` in the root of this repo.

## Configuring your application

To start an app with Devon, you will need a config file called `devon.conf.yaml` in that app's Git repo. It should look something like this:

```yaml
modes:
  development:
    command: ["bundle", "exec", "rails", "server"]
    dependencies:
      # These are key-value pairs, where the key is the name of
      # the dependency's git repo, and the value is the name of
      # the mode the dependency should be started in.
      my-app: dependency
      your-app: custom-mode
  dependency:
    command: ["dev/up.sh"]
    dependencies:
      my-app: dependency
      your-app: dependency
```


## Influencing ideas

### Apps should document how to run them in dev

This is an extension of the idea behind the `dev/up.sh` convention we've had for a long time now. Devon attempts to put a common interface on starting apps, so that if I'm hacking on an app for the first time, I don't have to worry about how to start it.

### Apps should document their dependencies on other apps

Just as we document our build-time dependencies with tools like `bundler` or `npm`, we should document our runtime dependencies. These are just as important when it comes to running the application.

### The best documentation can be used for automation, too!

Documentation can give us certainty. Once we have certainty (and consistent formatting), we can automate.

### Convention over (but not instead of!) configuration

Devon's conventions currently take the form of _stated assumptions_, as documented below. The aim is to gradually introduce some flexibility around those assumptions (e.g. by making them configurable), but keep defaults that make sense, so you don't have to configure *every little thing* in order to get started.

## Important assumptions

### All app repos live in the same place

By default, that place is `~/src`. You can override that by creating a config file called `.devon.yaml` in your home directory, with contents like this:

```yaml
source-code-base-dir: /path/to/your/code
```

### One app should never be started in 2 modes at the same time

I (Lucas) can't think of a sane use case for that, so I'm pretty sure this assumption is safe.

### Applications run in the background

Devon starts applications by running commands somewhat naively, waiting for each command to finish before starting the next. This means that applications can interrupt the starting process if they run in the foreground, instead of doing what we expect and starting a daemonised process. Because they never exit, subsequent applications wouldn't be able to start.

### All dependency relationships are the same

At the moment, a dependency relationship means "A depends on B to work properly, but A can start even if B isn't running." This is optimistic, since there is nothing stopping A from crashing if B is not available at startup.

This probably needs fixing with an actual dependency graph implementation -- the current implementation just builds a list that includes all the dependencies, with no information about which apps depend on which others. Hence, there's no reliable way to figure out what to start first.

## Development

**[Trello board](https://trello.com/b/MsxE9Nw6/devon-the-dev-application-starter)**

### Pre-commit hooks (optional)

This repo is set up to use `pre-commit` to share Git pre-commit hooks. The config is in `.pre-commit-config.yaml`. It's currently set up to check a few essential things:

* Files end with \n
* No trailing whitespace
* YAML formatting
* Golang formatting

This is to help prevent CI failures when Lucas forgets to run `go fmt` before pushing his commits. You can use it too, but you don't have to. Install with `brew install pre-commit`, then run `pre-commit install` in the root of this repo.

See https://pre-commit.com for more info.
