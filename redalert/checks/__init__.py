"""Contains various checks that redalert can perform on a system."""

from .exc import CheckFailure
from .file_exists import FileExistsCheck, FileNotExistsCheck
from .ulimit_checks import AddressSizeCheck, UlimitCheck


def get_check(name, args=None):
    '''Return the appropriate check instance based on test name.'''
    if args is None:
        args = {}

    if name == 'address-size':
        return AddressSizeCheck(**args)
    elif name == 'file-exists':
        return FileExistsCheck(**args)
    elif name == 'file-does-not-exist':
        return FileNotExistsCheck(**args)
    elif name == 'dpkg-installed':
        from .dpkg_installed import DpkgCheck
        return DpkgCheck(**args)
    elif name == 'compile-gcc':
        from .compile_gcc import CompileGccCheck
        return CompileGccCheck(**args)
    elif name in ['run-bash-script', 'run-script']:
        return RunScript(**args)

    raise NotImplementedError
