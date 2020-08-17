class CircularDependencyException < StandardError
  def initialize(apps)
    @apps = apps
    super("Circular dependency when resolving dependencies: #{apps_s}")
  end

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
end
