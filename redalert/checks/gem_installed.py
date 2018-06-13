"""Checks related to Ruby Gems."""

import subprocess

from .exc import CheckFailure


def get_installed_gems():
    """Return a list of installed gems"""
    proc = subprocess.run(['gem', 'list'])
    return str(proc.stdout).split('\n')


class GemInstalledCheck:
    """Checks if gems are installed and optionally checks version"""

    def __init__(self, name, version=""):
        self.name = name
        self.version = version

    def check(self):
        """Check if the gem indicated by name is installed and the correct version."""
        installed_gems = get_installed_gems()
        for gem in installed_gems:
            if self.name in gem:
                if self.version not in gem:
                    raise CheckFailure(
                        'gem {} installed but is version {} expected {}'\
                        .format(self.name, gem, self.version))
                return
        raise CheckFailure('{}')
