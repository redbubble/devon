#!/usr/bin/env ruby

## TODO: Figure out how to make this work regardless of the "current" version of
## Ruby according to rbenv, chruby or whatever. Options include:
#
# * Shebanging the system ruby directly
# * Writing ridiculously compatible code
# * Using another language, preferably a compiled one
#

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
  # TODO: Maybe a `stop` method is a good idea?
  #
  # But how do we make it work with dependencies? Do we just stop everything?
  # And if not, how do we decide what to stop, and what not to?
  # How do we even know what's been started with devon?
  #
  # AAAAAAAAAAAARGH!!! :(
  #

  def start(app, options)
    mode = options[:mode]

    # TODO: Do we want to ensure this app is up-to-date?
    #
    # * Clone it if absent?
    # * Pull it?
    # * Pull it if on master?
    # * Depends on the mode (e.g. not in development, yes in dependency?)
    # * Let the mode config specify?
    #

    config_path = File.join(SOURCE_CODE_BASE, app, CONFIG_FILE_NAME)
    config =
      begin
        YAML.load(File.read(config_path))
      rescue Errno::ENOENT
        # TODO: Decide what to do when we can't find a config file.
        #
        # Options:
        #
        # * Leave it and move on (current)
        # * Stop everything
        # * Just exit (e.g. by raising an unhandled error)
        # * Take some sensible default action (e.g. run ./dev/up.sh if in dependency mode)
        #
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
