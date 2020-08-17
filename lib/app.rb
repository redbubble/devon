require 'yaml'

class App

  class CouldNotReadConfigException < StandardError
    def initialize(app)
      @app = app
      super("#{app} has a problem with the YAML syntax in its #{App::CONFIG_FILE_NAME}, so devon can't start it :(")
    end
  end

  class NoConfigFileException < StandardError
    def initialize(app)
      @app = app
      super("#{app} doesn't have a #{App::CONFIG_FILE_NAME}, so devon can't start it :(")
    end
  end

  class ModeDoesNotExistException < StandardError
    def initialize(app, mode)
      @app = app
      @mode = mode
      super("#{app} doesn't have a mode called '#{mode}', so devon can't start it :(")
    end
  end

  SOURCE_CODE_BASE = "#{ENV['HOME']}/src"
  CONFIG_FILE_NAME = 'devon.conf.yaml'

  def initialize(name, mode)
    @name = name
    @mode = mode
    @config = read_config
  end

  attr_accessor :name, :mode, :config

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

  private

  def read_config
    config_path = File.join(SOURCE_CODE_BASE, name, CONFIG_FILE_NAME)
    YAML.load(File.read(config_path))
  rescue Errno::ENOENT
    raise NoConfigFileException.new(name)
  rescue Psych::SyntaxError
    raise CouldNotReadConfigException.new(name)
  end
end
