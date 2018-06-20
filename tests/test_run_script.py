import pytest
from redalert.checks import CheckFailure
from redalert.checks.run_script import RunScript


def test_run_bash_script_doesnt_raise_on_0():
    bs = RunScript(source='exit 0')
    bs.check()

def test_run_bash_script_raises_on_1():
    bs = RunScript(source='exit 1')
    with pytest.raises(CheckFailure):
        bs.check()

def test_run_bash_script_raises_on_bad_command():
    bs = RunScript(source='ls -l /fake_folder')
    with pytest.raises(CheckFailure):
        bs.check()

def test_run_bash_script_doesnt_raise_on_match_output():
    bs = RunScript(source='echo "123"', output="123")
    bs.check()

def test_run_bash_script_raises_on_mismatch_output():
    bs = RunScript(source='echo "111"', output="123")
    with pytest.raises(CheckFailure):
        bs.check()

def test_run_bash_script_doesnt_raise_on_good_interpreter():
    bs = RunScript(source='echo "111"', interpreter='/bin/sh')
    bs.check()

def test_run_bash_script_raises_on_bad_interpreter():
    bs = RunScript(source='echo "111"', interpreter='/bin/shoe')
    with pytest.raises(FileNotFoundError):
        bs.check()
