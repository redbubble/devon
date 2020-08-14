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

class App

  class CouldNotReadConfigException < StandardError
    def initialize(app)
      @app = app
      msg = "Aw, shoot! It seems as though #{app} has a problem with the YAML syntax in its #{AppStarter::CONFIG_FILE_NAME}, so devon can't start it :("
      super(msg)
    end
  end

  class NoConfigFileException < StandardError
    def initialize(app)
      @app = app
      msg = "Oh noes! It looks like #{app} doesn't have a #{AppStarter::CONFIG_FILE_NAME}, so devon can't start it :("
      super(msg)
    end
  end

  class ModeDoesNotExistException < StandardError
    def initialize(app, mode)
      @app = app
      @mode = mode
      msg = "Gadzooks! It appears that #{app} doesn't have a mode called '#{mode}', so devon can't start it :("
      super(msg)
    end
  end

  SOURCE_CODE_BASE = "#{ENV['HOME']}/src"
  CONFIG_FILE_NAME = 'devon.conf.yaml'

  def initialize(name, mode)
    @name = name
    @mode = mode
  end

  def start

    puts "Starting #{name} in #{mode} mode..."

    if Options.verbose?
      puts config
    end


    # Is it a bird?
    # Is it a plane?
    # No! It's depth-first recursion!
    #
    # TODO: Maybe handle some errors?
    #
    mode_config['dependencies'].each do |dep, dep_mode|
      new(dep, dep_mode).start
    end

    # Finally, actually start the thing
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

AppStarter.new.start(app, Options.mode)
