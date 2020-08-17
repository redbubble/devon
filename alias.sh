devon () {
    devon_repo=${HOME}/src/devon
    ruby_version=$(cat ${devon_repo}/.ruby-version)

    # Set the Ruby version in a way that rbenv will pay attention to.
    export RBENV_VERSION=${ruby_version}

    if command -v chruby-exec &> /dev/null; then
        prefix="chruby-exec ${ruby_version} --"
    fi

    ${prefix} ${devon_repo}/devon.rb
}
