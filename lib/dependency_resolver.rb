class DependencyResolver
  def initialize
    @apps = []
  end

  def add(app, dep_chain: [])
    # Check if it's possible to add the incoming app to our collection. Raise an
    # appropriate error if not.
    if dep_chain.include?(app.name)
      raise CircularDependencyException.new(dep_chain << app.name)
    end

    if @apps.include?(app.name)
      similar_app = apps.find { |a| a.name == app.name }

      if similar_app.mode != app.mode
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
    end

    # Add this app to our collection
    apps << app

    # Is it a bird?
    # Is it a plane?
    # No! It's depth-first recursion!
    app.dependencies.each do |dep, dep_mode|
      add(App.new(dep, dep_mode), dep_chain: dep_chain << app.name)
    rescue App::NoConfigFileException => ex
      puts "WARNING: #{ ex.message }"
      next
    end

    puts "Devon will start these apps:"
    apps.each { |a| puts "* #{a.name} (in '#{a.mode}' mode)" }
  end

  def start
    # Reverse order so that dependencies are started before their dependents.
    # This is not *guaranteed* to start dependencies before dependents,
    # especially in cases where a dependency has multiple dependents.
    #
    # If we're going to guarantee that, we'll need to implement a proper graph
    # in memory, and then resolve which things have dependencies, and which ones
    # we've started. Hopefully we can get away without it.
    #
    apps.reverse.each { |app| app.start }
  end

  private

  attr_reader :apps
end
