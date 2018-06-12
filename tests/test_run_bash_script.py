import pytest
from redalert.checks import CheckFailure
from redalert.checks.run_bash_script import RunBashScript


def test_run_bash_script_doesnt_raise_on_0():
    bs = RunBashScript(source='exit 0')
    bs.check()

def test_run_bash_script_raises_on_1():
    bs = RunBashScript(source='exit 1')
    with pytest.raises(CheckFailure):
        bs.check()

def test_run_bash_script_doesnt_raise_on_match_output():
    bs = RunBashScript(source='echo "123"', output="123")
    bs.check()

def test_run_bash_script_raises_on_mismatch_output():
    bs = RunBashScript(source='echo "111"', output="123")
    with pytest.raises(CheckFailure):
        bs.check()

def test_run_bash_script_doesnt_raise_on_good_interpreter():
    bs = RunBashScript(source='echo "111"', interpreter='/bin/sh')
    bs.check()

def test_run_bash_script_raises_on_bad_interpreter():
    bs = RunBashScript(source='echo "111"', interpreter='/bin/shoe')
    with pytest.raises(FileNotFoundError):
        bs.check()