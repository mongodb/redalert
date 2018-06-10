"""dpkg package checkers"""

import apt

from .exc import CheckFailure


class DpkgCheck:
    """Check if a package is installed in dpkg"""

    def __init__(self, name, version=None):
        self.name = name
        self.version = version

    def check(self):
        """Check that the package is installed and optionally the correct version."""
        cache = apt.Cache()
        pkg = cache[self.name]
        if not pkg.is_installed:
            raise CheckFailure('package {} is not installed'.format(self.name))
        if self.version is not None and pkg.installed.version != self.version:
            raise CheckFailure('package {} version {} is installed not {}'\
                            .format(self.name, pkg.installed.version, self.version))
