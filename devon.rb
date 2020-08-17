#!/usr/bin/env ruby

$LOAD_PATH << File.join(File.dirname(__FILE__), "lib")

require 'options'
require 'exceptions'
require 'dependency_resolver'
require 'app'

Options.parse!
puts Options.all

# If no app name is given, default to the name of the current directory
app = ARGV.empty? ? File.basename(ENV['PWD']) : ARGV.first

puts "Starting #{app} in '#{Options.mode}' mode"

resolver = DependencyResolver.new
app = App.new(app, Options.mode)

begin
  resolver.add(app)
rescue CircularDependencyException => ex
  puts "ERROR: Circular dependency when resolving dependencies for #{app.name}: #{ex.apps_s}"
  exit 1
end
