# redalert
A system image validation tool

## Installation / Usage

To install use pip:

    $ pip install redalert

Or clone the repo:

    $ git clone https://github.com/mongodb/redalert.git
    $ python setup.py install

Additionally standalone binaries are provided on the release page.

# Contributing

## Testing

To test redalert you need the following prerequisites:

    - vagrant
    - Virtualbox (for vagrant)
    - make

Then simply run 

```
make test
```

This will run all of the tests including "local"
tests on your machine. It will then spin up 3 vagrant VMs for testing the
platform specific code.

To only run the local or "platform agnostic" tests then run:

```
make unit_tests
```

To run a specific platform in vagrant one of:

```
make test_ubuntu
make test_centos
make test_windows
```

You can always invoke pytest yourself of course if you're on one of those
platforms and want to not run VMs with vagrant.

# Example

TBD
