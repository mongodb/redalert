# Design: Redalert Image Testing
## Author: Mathew Robinson <chasinglogic@gmail.com>

## Summary

Redalert is a tool used for system validation, similar to ChefSpec/ServerSpec.
It is written in Go and configured with YAML and should be familiar to those
accustomed to working with Ansible.

## Behavioral Description

Redalert will read yaml files containing tests and test configuration, it will
then run these tests against the system being tested and report the successes
and failures.

## Detailed Design

### CLI

Redalert will have a subcommand interface with the following commands:

 - [init](#init)
 - [run](#run)
 - [convert](#convert)
 - [alias](#alias)
 - [gen-aliases](#gen-aliases)

#### init

The init subcommand will take zero to one arguments. The argument will be a
directory name in which to create the dummy test files. If no argument is given
Redalert will generate them in the local directory.

The test files generated will be:

 - tests.(yml|yaml)
 - aliases.(yml|yaml)

If git is detected in $PATH then init will attempt to also run git init.

#### run

Run will simply run the tests. It takes zero to one arguments. If no arguments
are given it looks for a tests.(yml|yaml) or redalert.(yml|yaml) in the local
directory, in `$HOME/.redalert/` (`%APPDATA%` instead of `$HOME` on windows.)
and finally in `/etc/redalert/` (`C:\redalert` on windows). 

It takes the following flags:

 - `--quiet` only report test failures.
 - `--jobs $VALUE` specify the number of parallel tests to run. Default to # of cores.
 - `--test $TEST_NAME` specify a test by name to run, can be provided multiple times.
 - `--suite $SUITE_NAME` which test suite to run, an [alias](#aliases) can be provided as the suite name.
   Can be provided multiple times. There is a default "all" alias which matches
   all tests and is used if `--suite` is not provided. 
 - `--file $FILE_PATH` in lieu of searching in directories you can specify a configuration yaml file using this flag.
   If provided, the aliases.(yml|yaml) file will not be looked for but aliases in the 
   given file will still be loaded.
 - `--output $FORMAT` specify the output format. Valid values: text, json, csv. Default: text.

Run will attempt to load an additional file `aliases.(yml|yaml)` which
specifies the available [aliases](#aliases). If not found, aliases will be
looked for as a key in the `tests.(yml|yaml)` according to the
[aliases](#aliases) section.

Run will run tests in parallel.

#### convert 

Convert converts other server testing formats to redalert. This will only
support [Greenbay](https://github.com/mongodb/greenbay) to start however will
support other common formats in the future such as ServerSpec. 

#### alias 

alias takes two arguments followed by varargs of the form:

``` 
redalert alias TEST_FILE ALIAS_NAME SUITES_TO_ADD_TO_ALIAS...
```

Redalert will attempt to find an aliases.(yml|yaml) in the same directory as
`TEST_FILE`. If found, it will load both files and replace the
aliases.(yml|yaml) with a new aliases file which includes the existing and new
alias. Otherwise, it will create an alias in aliases.(yml|yaml) it will then
populate the alias with the `SUITES_TO_ADD_TO_ALIAS` suites replacing in any
tests, which have all those suites, the suites with the create alias. It will
then save the test file.

#### gen-aliases (Stretch Goal)

gen-aliases takes a test file as input and generates an aliases.(yml|yaml) with
aliases that are generated based on which suites are commonly used together.
This command is computationally intensive so likely will take a long time.

### Config File/s

Redalert has two config files: aliases.(yaml|yml) and tests.(yaml|yml). They
are described in the below sections and should always be co-located on the file
system.

#### Aliases

Aliases are a way of uncluttering your test files. For example if you are
testing various versions of Ubuntu and you have suites names
ubuntu$VERSION\_SPECIFIER, such that your list of ubuntu suites looks like:

 - `ubuntu1204`
 - `ubuntu1404`
 - `ubuntu1604`
 - `ubuntu1804`

And you have a test called "gcc installed" that looks like:

```yaml
tests:
  - name: gcc installed
    suites:
      - ubuntu1204
      - ubuntu1404
      - ubuntu1604
      - ubuntu1804
    type: dpkg-installed
    args:
      name: gcc
```

You can simplify this with an alias:

```yaml
aliases:
    ubuntu:
      - ubuntu1204
      - ubuntu1404
      - ubuntu1604
      - ubuntu1804
tests:
  - name: gcc installed
    suites:
      - ubuntu
    type: dpkg-installed
    args:
      name: gcc
```

Which for this trivial example is not particularly useful. But as you
accumulate multiple tests which span all Ubuntus you can see that your test
file would grow quite long quite quickly. Additionally if you only have a set
of suites which run every test it makes removing a deprecated platform a one
line change.

#### Tests

The tests yaml file has a single key "tests" which is a yaml list of maps which have the following form:

```yaml
# Human readable test name
name: name of the test
# Suites are user defined, this key is optional if not provided test will only
# match the all alias
suites:
    - suite this test should be run for
# See (Checks "tests") below
type: test-type
# Arguments for the test. See the documentation for the given test type to
# determine which args are available / required.
args:
    key: value
```

The final product for a yum-installed test looks like:

```yaml
tests:
  - name: gcc is installed
    suites:
      - rhel70
    type: yum-installed
    args:
      name: gcc
```

### Checks ("tests")

The Checker interface is shown below:

```go
// Args is a convenience type to express the "args" key in the test yaml
type Args map[string]interface{}

type Checker interface {
  Check() error
}

// Argable is any struct which can create a Checker from the YAML args we get
// back from a test block.
type Argable interface {
	FromArgs(args Args) (Checker, error)
}
```

Checks (synonym for tests in redalert terms) are structs which implement the
Check interface and follow a standard godoc comment template as laid out below:

```go
// ExampleCheck is just an example so it checks nothing.
// 
// Type: example-check
// 
// Supported Platforms:
//     - Mac
//     - Linux
//     - Windows
// 
// Requirements:
//     This test has no requirements. I could omit this section of the
//     docstring since it will work anywhere redalert can run. But if it
//     required that certain programs be available it should be documented
//     here.
// 
// Arguments:
//     required_arg (required): A string value that will be thrown away.
//     optional_arg: A value of any type that will be ignored. As a style
//                   guide example, I have inflated the description of this
//                   useless argument arbitrarily. This allows me to
//                   demonstrate how to wrap lines for long argument
//                   descriptions.
// 
// Notes:
//     If I had any special behavior based on arguments passed I should
//     note it here. For example if I have multiple Types and have 
//     different behavior for each I should make note of that here.
//     Otherwise this section can be omitted.
type ExampleCheck struct {
    RequiredArg string `mapstructure:"required_arg"`
    OptionalArg interface{} `mapstructure:"optional_arg"`
}

func (ec ExampleCheck) FromArgs(args Args) (Check, error) {
    if err := requiredArgs(args, "required_arg"); err != nil {
        return nil, err
    }
    
    err := decodeFromArgs(args, &ec)
    return ec, err
}

func (ec ExampleCheck) Check() error {
    return nil
}
```

The doc comments for the check must follow the format perscribed above as it is
used in the generated documentation for the checks.

You can specify multiple Types in a list form similiar to Supported Platforms.

All check classes must have a check method which takes no arguments and returns
no value. It should return an error if the test is a failure or if any setup
errors occur with the string error message to be shown to the user.

##### Adding a Checker

The process for writing a checker goes as follows:

 - Decide on check functionality
 - Decide on check "type name" (as seen in the yaml)
 - Write a check struct that implements the [Checker](https://github.com/mongodb/redalert/blob/master/checks/checks.go#L20) interface
 - Implement [Argable](https://github.com/mongodb/redalert/blob/master/checks/checks.go#L15) for your struct 
   - Often you can just copy the FromArgs of ExampleCheck [above](#checks-tests)
 - See the [file checker](https://github.com/mongodb/redalert/blob/master/checks/file.go) for a good example to follow
 - Finally, add your check to [load.go](https://github.com/mongodb/redalert/blob/master/checks/load.go) in the availableChecks map
   - Again see the "file-exists" and "file-does-not-exist" in the availableChecks map for good examples

#### MVP Check Types

These are the check types we will minimally need to implement to make Redalert
a drop-in Greenbay replacement for Build Team use.

 - [ ] type: address-size
 - [ ] type: command-group-all SUPERSEDED by run-bash-script
 - [ ] type: compile-and-run-gcc-auto SUPERSEDED by compile-gcc with arg run
 - [ ] type: compile-gcc-auto SUPERSEDED by compile-gcc
 - [ ] type: dpkg-installed
 - [ ] type: file-does-not-exist
 - [ ] type: file-exists
 - [ ] type: gem-installed
 - [ ] type: open-files
 - [ ] type: run-bash-script
 - [ ] type: shell-operation SUPERSEDED by run-bash-script
 - [ ] type: python-module-version
 - [ ] type: run-program-system-python (should be superseded by run-python-script)
 - [ ] type: run-program-system-python2 (should be superseded by run-python-script)
 - [ ] type: yum-group-any
 - [ ] type: yum-installed
 - [ ] type: compile-visual-studio
 - [ ] type: irp-stack-size
 - [ ] type: lxc-containers-configured (possibly not needed?)

## Future Work

 - Ansible Module / Integration
 - ServerSpec test parity
