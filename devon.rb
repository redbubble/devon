#!/usr/bin/env ruby

require 'optparse'
require 'yaml'

SOURCE_CODE_BASE = "#{ENV['HOME']}/src"
CONFIG_FILE_NAME = 'devon.conf.yaml'

# Initialise with some default values
@options = {
  mode: 'dependency'
}

# Parse the command line arguments
OptionParser.new do |opts|
  opts.on('-m', '--mode MODE', "The mode to run in, e.g. 'development' or 'dependency'") do |mode|
    @options[:mode] = mode
  end

  opts.on('-v', '--verbose', "Print all the informations!") do
    @options[:verbose] = true
  end
end.parse!

class AppStarter
  def start(app, options)
    mode = options[:mode]

    config_path = File.join(SOURCE_CODE_BASE, app, CONFIG_FILE_NAME)
    config =
      begin
        YAML.load(File.read(config_path))
      rescue Errno::ENOENT
        puts "Oh noes! It looks like #{app} doesn't have a #{CONFIG_FILE_NAME}, so devon can't start it :("
        return
      end

    puts "Starting #{app} in #{mode} mode..."

    if options[:verbose]
      puts config
    end

    raise "Could not find a mode named '#{mode}' for application '#{app}'." unless config['modes'].has_key?(mode)

    mode_config = config['modes'][mode]

    # Is it a bird?
    # Is it a plane?
    # No! It's depth-first recursion!
    #
    # TODO: Maybe handle some errors?
    #
    mode_config['dependencies'].each do |dep, dep_mode|
      start(dep, options.merge({mode: dep_mode}))
    end

    # Finally, actually start the thing
    # This is naive af, but it should be OK for a PoC...
    command = mode_config['command'].join(" ")

    if options[:verbose]
      puts "Running command: '#{command}'"
    end

    system("cd #{File.join(SOURCE_CODE_BASE, app)}; #{command}")
  end
end

puts @options

# If no app name is given, default to the name of the current directory
app = ARGV.empty? ? File.basename(ENV['PWD']) : ARGV.first

AppStarter.new.start(app, @options)
