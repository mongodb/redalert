import pytest
from redalert.checks import get_check


@pytest.mark.skippable
def test_supported_test_types():
    supported_test_types = [
        'address-size', 'compile-gcc', 'compile-visual-studio',
        'dpkg-installed', 'file-does-not-exist', 'file-exists',
        'gem-installed', 'irp-stack-size', 'lxc-containers-configured',
        'open-files', 'python-module-version', 'run-bash-script',
        'run-program-system-python', 'run-program-system-python2',
        'shell-operation', 'yum-group-any', 'yum-installed'
    ]

    for supported_test in supported_test_types:
        try:
            get_check(supported_test, args={})
        except NotImplementedError:
            pytest.fail('{} is not implemented'.format(supported_test))
        except TypeError:
            pass
        except ImportError as e:
            # Ignore apt on non-ubuntu platforms
            if 'apt' in str(e):
                pass
            else:
                raise e
