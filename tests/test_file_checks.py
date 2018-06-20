import pytest
from redalert.checks import CheckFailure
from redalert.checks.file_exists import FileExistsCheck, FileNotExistsCheck


@pytest.mark.agnostic
def test_file_exists(tmpdir):
    tmp = tmpdir.mkdir('file_exists').join('file.txt')
    tmp.write('text')
    fc = FileExistsCheck(name=tmp)
    fc.check()


@pytest.mark.agnostic
def test_file_exists_throws_on_no_exist():
    fc = FileExistsCheck(name='NOT_A_FILE')
    with pytest.raises(CheckFailure):
        fc.check()


@pytest.mark.agnostic
def test_file_exists_reverse():
    fc = FileNotExistsCheck(name='NOT_A_FILE')
    fc.check()


@pytest.mark.agnostic
def test_file_exists_reverse_throws_on_exist(tmpdir):
    tmp = tmpdir.mkdir('file_exists_reverse').join('file.txt')
    tmp.write('text')
    fc = FileNotExistsCheck(name=tmp)
    with pytest.raises(CheckFailure):
        fc.check()
