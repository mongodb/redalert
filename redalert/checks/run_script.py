# type: run-script

import subprocess
import tempfile

from .exc import CheckFailure

class RunScript:
    """A check for running a bash script and optionally checks output.

    Takes the following arguments:
    source (required): The source code of the script.
    output: string to which the output of the script will be compared to determine a successful run. If omitted, only checks that the script exits with returncode 0.
    interpreter: path to bash. Default is '/bin/bash'.
    """
    
    def __init__(self, source, output="", interpreter="/bin/bash"):
        self.source = source
        self.output = output
        self.interpreter = interpreter
        
    def check(self):
        script = tempfile.NamedTemporaryFile(mode='w+')
        script.write(self.source)
        script.flush()

        execution = subprocess.run([self.interpreter, script.name],
                                  stdout=subprocess.PIPE)
        
        if execution.returncode != 0:
            raise CheckFailure('Failed to execute the script. Return code: {}'.format(execution.returncode))
            
        exec_output = execution.stdout.decode('utf-8').strip() if execution.stdout else ''
        if self.output and exec_output != self.output:
            raise CheckFailure('Output {} doesn\'t match expected {}'.format(exec_output, self.output))
