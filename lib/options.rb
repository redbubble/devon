require 'optparse'

class Options
  @options = {
    mode: 'development'
  }

  def self.parse!
    # Parse the command line arguments
    OptionParser.new do |opts|
      opts.on('-m', '--mode MODE', "The mode to run in, e.g. 'development' or 'dependency'. Default: development") do |mode|
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
