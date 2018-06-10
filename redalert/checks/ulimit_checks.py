"""Ulimit checks"""

import resource

from .exc import CheckFailure


class UlimitCheck:
    def __init__(self, value):
        self.value = value

    def check(self):
        _soft, hard = resource.getrlimit(self.ulimit)
        if hard < self.value:
            raise CheckFailure('{} limit is {} which is less than {}'\
                            .format(self.ulimit, hard, self.value))


class AddressSizeCheck(UlimitCheck):
    ulimit = resource.RLIMIT_AS
