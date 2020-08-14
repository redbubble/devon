#!/usr/bin/env ruby

require 'optparse'
require 'yaml'

class Options
  @options = {
    mode: 'dependency'
  }

  def self.parse!
    # Parse the command line arguments
    OptionParser.new do |opts|
      opts.on('-m', '--mode MODE', "The mode to run in, e.g. 'development' or 'dependency'") do |mode|
        @options[:mode] = mode
      end

      opts.on('-v', '--verbose', "Print all the informations!") do
        @options[:verbose] = true
      end
    end.parse!
  end

  def self.all
    @options
  end

  def self.mode
    @options[:mode]
  end

  def self.verbose?
    @options[:verbose]
  end
end

class CircularDependencyException < StandardError
  def initialize(apps)
    @apps = apps
    super("Circular dependency when resolving dependencies: #{apps_s}")
  end

  private

  def apps_s
    @apps.join('->')
  end
end

class DependencyConflictException < StandardError
  def initialize(new_app, existing_app)
    @new_app = new_app
    @existing_app = existing_app
    super("Dependency conflict when resolving dependencies: #{new_app.name} in '#{new_app.mode}' mode conflicts with #{existing_app.name} in '#{existing_app.mode}' mode.")
  end

  private

  def apps_s
    @apps.join('->')
  end
end

class DependencyResolver
  def initialize()
    @apps = []
  end

  def add(app, dep_chain: [])
    # Check if it's possible to add the incoming app to our collection. Raise an
    # appropriate error if not.
    if dep_chain.has_key?(app.name)
      raise CircularDependencyException.new(dep_chain)
    end

    if @apps.has_key?(app.name)
      if apps[app.name].mode != app.mode
        # We're trying to depend on the same app in two different modes. This has
        # so much potential to go wrong that we're better off not even trying.
        raise DependencyConflictException.new(app, apps[app.name])
      else
        # We already have this dependency on our list, so there's nothing to do.
        #
        # But wait, isn't this a circular dependency? Well, no. We checked for
        # that above. We can get here when A depends on B and C, and both B and
        # C depend on D, which is allowed.
        #
        return
    end

    # Add this app to our collection
    apps << app

    # Is it a bird?
    # Is it a plane?
    # No! It's depth-first recursion!
    app.dependencies.each do |dep, dep_mode|
      add(App.new(dep, dep_mode), dep_chain << app.name)
    end
  end

  def start
    # Reverse order so that dependencies are started before their dependents.
    apps.reverse.each { |app| app.start }
  end

  private

  attr_reader :apps
end

class App

  class CouldNotReadConfigException < StandardError
    def initialize(app)
      @app = app
      super("Aw, shoot! It seems as though #{app} has a problem with the YAML syntax in its #{AppStarter::CONFIG_FILE_NAME}, so devon can't start it :(")
    end
  end

  class NoConfigFileException < StandardError
    def initialize(app)
      @app = app
      super("Oh noes! It looks like #{app} doesn't have a #{AppStarter::CONFIG_FILE_NAME}, so devon can't start it :(")
    end
  end

  class ModeDoesNotExistException < StandardError
    def initialize(app, mode)
      @app = app
      @mode = mode
      super("Gadzooks! It appears that #{app} doesn't have a mode called '#{mode}', so devon can't start it :(")
    end
  end

  SOURCE_CODE_BASE = "#{ENV['HOME']}/src"
  CONFIG_FILE_NAME = 'devon.conf.yaml'

  def initialize(name, mode)
    @name = name
    @mode = mode
  end

  attr_accessor :app, :mode

  def dependencies
    mode_config['dependencies']
  end

  def start

    puts "Starting #{name} in #{mode} mode..."

    if Options.verbose?
      puts config
    end

    if Options.verbose?
      puts "Running command: '#{command}'"
    end

    system("cd #{File.join(SOURCE_CODE_BASE, name)}; #{command}")
  end

  # This is naive af, but it should be OK for a PoC...
  def command
    mode_config['command'].join(" ")
  end

  def mode_config
    unless config['modes'].has_key?(mode)
      raise ModeDoesNotExistException.new(name, mode)
    end

    config['modes'][mode]
  end

  def config
    return @config if @config

    config_path = File.join(SOURCE_CODE_BASE, name, CONFIG_FILE_NAME)
    @config =
      begin
        YAML.load(File.read(config_path))
      rescue Errno::ENOENT
        raise NoConfigFileException.new(name)
      rescue Psych::SyntaxError
        raise CouldNotReadConfigException.new(name)
      end
  end
end

Options.parse!
puts Options.all

# If no app name is given, default to the name of the current directory
app = ARGV.empty? ? File.basename(ENV['PWD']) : ARGV.first

App.new.start(app, Options.mode)
