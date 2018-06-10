"""GCC compile checks"""

import os
import shlex
import shutil
import subprocess
import tempfile

from .exc import CheckFailure


class CompileGccCheck:
    """A check for compiling and optionally running c source code.

    Takes the following arguments:
    source (required): The source code of the C program.
    cflags: Additional flags to be passed to gcc. Can be a string or list.
    cflags_command: A shell command whose output will be used as flags to GCC

    """

    def __init__(self,
                 source,
                 cflags=None,
                 cflags_command="",
                 gcc=None,
                 run=False,
                 output=""):
        if cflags is None:
            cflags = []
        elif cflags is list:
            cflags = cflags
        elif cflags is str:
            cflags = cflags.split(" ")

        if gcc is None:
            self.gcc = 'gcc'

        self.source = source
        self.cflags = cflags
        self.cflags_command = cflags_command
        self.run = run
        self.output = output

    def check(self):
        """Check the C program given the args"""
        tmpdir = tempfile.mkdtemp()
        source_file = os.path.join(tmpdir, 'check.c')
        bin_file = os.path.join(tmpdir, 'check.out')

        with open(source_file) as tmp:
            tmp.write(self.source)

        if self.cflags_command:
            proc = subprocess.run(
                shlex.split(self.cflags_command),
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE)
            self.cflags += str(proc.stdout).split(' ')

        cmd = [self.gcc, '-Werror', '-o', bin_file, '-c', source_file] +\
            self.cflags
        proc = subprocess.run(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if proc.returncode != 0:
            output = str(proc.stdout) + str(proc.stderr)
            raise CheckFailure('command {} failed with code {}: {}'\
                            .format(cmd, proc.returncode, output))
        if self.run:
            compiled_proc = subprocess.run([bin_file], stdout=subprocess.PIPE)
            if compiled_proc.returncode != 0:
                raise CheckFailure('compiled program exited with code {}: {}'\
                                .format(compiled_proc.returncode,
                                        str(compiled_proc.stdout)))

            if self.output and self.output != str(compiled_proc.stdout):
                raise CheckFailure(
                    'compiled program output did not match expected {}: {}'\
                    .format(self.output, str(compiled_proc.stdout)))

        # Only clean up if successful since the leftover state can be
        # helpful for debugging a failure.
        shutil.rmtree(tmpdir)
