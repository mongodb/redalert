"""Ulimit checks"""

import resource

from .exc import CheckFailure


class UlimitCheck:
    """Checks the ulimit specified by a classes ulimit static member.

    This class is meant to be subclassed and should not be used directly.
    """
    ulimit = None

    def __init__(self, value):
        self.value = value

    def check(self):
        """Check that this classes ulimit is greater than or equal to the given value."""

        if self.ulimit is None:
            raise CheckFailure('UlimitCheck called directly. No ulimit to check.')

        _soft, hard = resource.getrlimit(self.ulimit)
        if hard <= self.value:
            raise CheckFailure('{} limit is {} which is less than {}'\
                            .format(self.ulimit, hard, self.value))


class AddressSizeCheck(UlimitCheck):
    """Check that the ulimit Address Size is equal to or greater than the given
    value."""
    ulimit = resource.RLIMIT_AS
