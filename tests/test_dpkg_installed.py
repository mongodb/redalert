import pytest
from redalert.checks.exc import CheckFailure


@pytest.mark.ubuntu
def test_dpkg_installed():
    from redalert.checks.dpkg_installed import DpkgCheck

    check = DpkgCheck("python3")
    check.check()

    check = DpkgCheck("python3", version='3.6.5-3')
    check.check()

    check = DpkgCheck("not-an-ubuntu-package")
    with pytest.raises(CheckFailure):
        check.check()

    check = DpkgCheck("python3", version='2.7')
    with pytest.raises(CheckFailure):
        check.check()
