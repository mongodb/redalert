#!/usr/bin/env python3
"""This source codifies our deprecations and argument modifications to the
Greenbay.yaml format. It takes a single argument: the path to the greenbay.yaml and will create a redalert.yaml from it."""

import sys
from collections import OrderedDict

import pyaml
import yaml


# Make it so yaml multi line strings use the | syntax that we prefer
def str_presenter(dumper, data):
    data = data.rstrip()
    if len(data.splitlines()) > 1:  # check for multiline string
        return dumper.represent_scalar(
            'tag:yaml.org,2002:str', data, style='|')
    if data.startswith('{') or data.endswith('}'):
        return dumper.represent_scalar(
            'tag:yaml.org,2002:str', data, style="'")
    return dumper.represent_scalar('tag:yaml.org,2002:str', data, style='')


pyaml.add_representer(str, str_presenter)


def represent_ordereddict(dumper, data):
    value = []

    if 'name' in data and 'suites' in data:
        for key in ['name', 'suites', 'type', 'args']:
            value.append((dumper.represent_data(key),
                          dumper.represent_data(data[key])))
    else:
        for item_key, item_value in data.items():
            node_key = dumper.represent_data(item_key)
            node_value = dumper.represent_data(item_value)
            value.append((node_key, node_value))

    return yaml.nodes.MappingNode(u'tag:yaml.org,2002:map', value)


pyaml.add_representer(OrderedDict, represent_ordereddict)
pyaml.add_representer(dict, represent_ordereddict)


def convert_shell_operation(test):
    """shell-operation superseded by run-bash-script"""
    cmd = test['args'].pop('command')
    test['args']['source'] = cmd.strip()
    test['type'] = 'run-bash-script'
    return test


def convert_run_program_system_python(test):
    """run-program-system-python* superseded by run-python-script"""
    if 'windows' in test['name'] and test['type'].endswith('python2'):
        test['args']['interpreter'] = 'C:\\Python27\\python.exe'
    elif 'windows' in test['name']:
        test['args']['interpreter'] = 'C:\\Python36\\python.exe'
    elif test['type'].endswith('python2'):
        test['args']['interpreter'] = 'python2'
    test['type'] = 'run-python-source'
    return test


def convert_command_group_all(test):
    """command-group-all superseded by run-bash-script and run-powershell-source"""
    cmds = [x['command'] for x in test['args'].pop('commands')]
    test['args']['source'] = '\n'.join(cmds)
    if 'windows' in test['name']:
        test['type'] = 'run-powershell-source'
    return test


def convert_compile_gcc(test):
    """compile-*gcc-auto superseded by compile-gcc"""
    if 'and-run' in test['type']:
        test['args']['run'] = True
    test['type'] = 'compile-gcc'
    return test


TRANSLATION_TABLE = {
    'shell-operation': convert_shell_operation,
    'run-program-system-python': convert_run_program_system_python,
    'run-program-system-python2': convert_run_program_system_python,
    'command-group-all': convert_command_group_all,
    'compile-and-run-gcc-auto': convert_compile_gcc,
    'compile-gcc-auto': convert_compile_gcc
}


def convert(test):
    """Convert a test if necessary."""
    if test['type'] in TRANSLATION_TABLE:
        print('Found superseded test:', test['type'])
        translation = TRANSLATION_TABLE[test['type']]
        print(translation.__doc__)
        return translation(test)
    return test


def main():
    greenbay_yaml_path = sys.argv[1]
    with open(greenbay_yaml_path) as gyml:
        greenbay_yaml = yaml.load(gyml)

    redalert_yaml = {'tests': []}
    tests = greenbay_yaml['tests']
    for test in tests:
        redalert_yaml['tests'].append(convert(test))

    with open('redalert.yaml', 'w') as ryml:
        ryml.write(pyaml.dump(redalert_yaml))


if __name__ == '__main__':
    main()
