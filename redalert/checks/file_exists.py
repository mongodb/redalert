"""File existence and absence checking."""

import os

from .exc import CheckFailure


class FileExistsCheck:
    """Checks if file at name exists"""

    def __init__(self, name, reverse=False):
        self.name = name
        self.reverse = reverse

    def check(self):
        """Check for the files existence"""
        if not os.path.isfile(self.name):
            raise CheckFailure('file {} not found'.format(self.name))


class FileNotExistsCheck:
    """Checks if file at name does not exist"""

    def __init__(self, name):
        self.name = name

    def check(self):
        """Check for the files absence"""
        if os.path.isfile(self.name):
            raise CheckFailure('file {} exists and should not'.format(
                self.name))
