#rpenv

displays env vars set from existing environment ( [skippable](#usage) ) and loaded from config file in specified environment (ci, qa, or prod) or executes command in the context of the existing environment variables and ones loaded from a config file.
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Usage](#usage)
- [Installation](#installation)
  - [macOS](#macos)
  - [Other OS](#other-os)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
## Usage

    $ rpenv -v

or

    $ rpenv -version

displays rpenv version information

    $ rpenv [-skip-local] <env>

or

    $ rpenv [-skip-local] <env> <cmd>

where `<env>` is one of `ci`, `qa`, or `prod` (`production` should also work) and `<cmd>` is the desired command you wish to run. If called without a `<cmd>`, `rpenv` will return a list of all the env vars in the `/etc/rentpath/environment.cfg` file merged with your current environment variables (i.e. whatever `/usr/bin/env` would return), or with `-skip-local`, your current environment variables will not be merged (the envs are still there from the parent process, just not displayed). When run with a `<cmd>`, it will execute that `<cmd>` after setting the environment with the values returned if `rpenv` is run without a `<cmd>`.

## Testing

For testing without hitting the network.

`go test -v`

For testing with hitting the network.

`go test -v -system`

## Installation

You will need to have `~/.config/.rpenv` with ini style configuration for `ci`, `qa`, and `prod` URIs
You can find those configurations in the [rpenv configuration wiki page](https://github.com/rentpath/idg/wiki/rpenv-configuration).

### macOS
    brew update && brew tap rentpath/homebrew && brew install rentpath/homebrew/rpenv

### Other OS
Install requirements for GO, to build the binary.

    yum -y install go

Clone the repo

    git clone git@github.com:rentpath/rpenv.git "$(go env GOPATH)/src/rpenv"

Build the go binary

    cd "$(go env GOPATH)/src/rpenv"
    go build

Set up build environment

    yum install build-tools
    mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

Copy files to build directories and build RPM

    mv rpenv ~/rpmbuild/SOURCES/
    cp rpenv.spec ~/rpmbuild/SPECS/
    rpm -bb ~/rpmbuild/SPECS/rpenv.spec

##License
[MIT](https://github.com/rentpath/rpenv/blob/master/LICENSE)
