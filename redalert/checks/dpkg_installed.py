"""dpkg package checkers"""

import subprocess

from .exc import CheckFailure


class DpkgCheck:
    """Check if a package is installed in dpkg"""

    def __init__(self, name, version=None):
        self.name = name
        self.version = version

    def check(self):
        """Check that the package is installed and optionally the correct version."""
        proc = subprocess.run(
            ['dpkg', '-s', self.name], stdout=subprocess.PIPE)
        if proc.returncode != 0:
            raise CheckFailure('package {} is not installed'.format(self.name))

        # End early if version check not necessary
        if self.version is None:
            return

        output = str(proc.stdout.decode('utf-8'))
        version = ''
        for line in output.splitlines():
            if line.startswith('Version:'):
                split = line.split(':')
                version = split[len(split) - 1]
                break

        # remove leading and trailing spaces
        version = version.strip()
        if version != self.version:
            raise CheckFailure('package {} version {} is installed not {}'\
                            .format(self.name, version, self.version))
