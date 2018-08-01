package checks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/google/shlex"
	"golang.org/x/sys/windows/registry"
)

func init() {
	availableChecks["compile-visual-studio"] = CompileVisualStudioFromArgs
}

// CompileVisualStudio runs VisualStudio compile.
//
// Type:
//	 - compile-visual-studio
//
// Support Platforms:
//   - Windows
//
// Arguments:
//   source (required): The source code to compile.
//   cflags: A string that will be parsed using shlex and passed as arguments to cl.exe
//   version: Visual studio version to use. Default is the latest version installed on the system
//   extension: The file extension to use for the generated temporary file. Default is "cpp"
//   run: If true try running the compiled binary
//   compiler: A string path to the compiler to use. This is a template that gets
// 	           one variable: '{{ .VisualStudioPath }}' which is the path to the
//             Visual Studio installation we are using, ending with the filepath separator
//             '\'. For example for VS2015 it would be:
//
// 	           C:\Program Files (x86)\Microsoft Visual Studio 14.0\
//
// 	           This allows you to specify compilers other than VC\bin\cl.exe to use.
// 	           For example if you wanted to compile C#. You can pass a hardcoded string here
// 	           without the template variable if you wish to specify a full path and
// 	           bypass our VisualStudio installation detection logic.
//
// 	           This defaults to "{{ .VisualStudioPath }}VC\\bin\\cl.exe"
//
// Notes:
//   For the version argument you can either specify the year (2015) or the actual
//   version number as recognized by Visual Studio registry entries (14.0) and
//   this check will do the math to figure out what you want.
type CompileVisualStudio struct {
	Source    string
	Cflags    string
	Version   float64
	Extension string
	Compiler  string
	Run       bool
}

func getRealVersionNumber(version float64) float64 {
	s := strconv.FormatFloat(version, 'f', -1, 64)
	// If the version is 4 digits long it's assumed to be a year
	if len(s) == 4 {
		year := int(version)
		// 2010 was version 10.0 then they waited 3 years to release version 12.0
		if year == 2010 {
			return 10.0
		}

		// 2013 and 2015 were the first "2 year cadence" releases and were
		// increased two version numbers per release
		if year <= 2015 {
			return 10.0 + float64(year-2011)
		}

		// After 2015 came 2017 (latest at the time of this writing) despite being
		// released 2 years after 2015 it was only bumped one version number (15.0)
		// this code assumes they'll be consistent and keep going up by one
		// every two years, when 2019 drops we may need to update this function
		// again.
		return 14.0 + float64((year-2015)/2)
	}

	return version
}

const visualStudioRegPath = "SOFTWARE\\Wow6432Node\\Microsoft\\VisualStudio\\SxS\\VS7"

func findVisualStudio(version float64) (string, error) {
	// Convert a year, if given, to a known visual studio version number
	realVersionNumber := getRealVersionNumber(version)
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE, visualStudioRegPath, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("unable to open registry: %s: %s", visualStudioRegPath, err)
	}

	defer key.Close()

	// If version specified just get it
	if version != 0.0 {
		versionString := strconv.FormatFloat(realVersionNumber, 'f', 1, 64)
		val, _, err := key.GetStringValue(versionString)
		return val, err
	}

	installedVersions, err := key.ReadValueNames(-1)
	if err != nil {
		return "", fmt.Errorf("Unable to read installed Visual Studio verisons: %s", err)
	}

	sort.Strings(installedVersions)

	val, _, err := key.GetStringValue(installedVersions[len(installedVersions)-1])
	return val, err
}

// Check compiles, and optionally runs, code using VisualStudio
func (cvs CompileVisualStudio) Check() error {
	vsPath, err := findVisualStudio(cvs.Version)
	if err != nil {
		return fmt.Errorf("Problem finding Visual Studio: %s", err)
	}

	tmpfolder, err := ioutil.TempDir("", "compileVisualStudio_")
	if err != nil {
		return fmt.Errorf("Problem creating a tmpdir: %s", err)
	}

	// defer os.RemoveAll(tmpfolder)

	if cvs.Extension == "" {
		cvs.Extension = "cpp"
	}

	srcfileName := filepath.Join(tmpfolder, "src."+cvs.Extension)
	outfileName := filepath.Join(tmpfolder, "out.exe")

	srcfile, err := os.Create(srcfileName)
	if err != nil {
		return fmt.Errorf("Problem creating a srcfile: %s", err)
	}

	cvs.Source = strings.Replace(cvs.Source, "\n", "\r\n", -1)
	content := []byte(cvs.Source)

	if _, err := srcfile.Write(content); err != nil {
		return fmt.Errorf("Problem writing to a tmpfile: %s", err)
	}

	if err := srcfile.Close(); err != nil {
		return fmt.Errorf("Problem closing the tmpfile: %s", err)
	}

	compileCommandTemplate := template.New("compileCommand")
	if cvs.Compiler == "" {
		compileCommandTemplate, err = compileCommandTemplate.Parse("{{ .VisualStudioPath }}VC\\bin\\cl.exe")
	} else {
		compileCommandTemplate, err = compileCommandTemplate.Parse(cvs.Compiler)
	}

	if err != nil {
		return fmt.Errorf("Problem parsing compileCommand template: %s", err)
	}

	buffer := &bytes.Buffer{}
	err = compileCommandTemplate.Execute(buffer, struct {
		VisualStudioPath string
	}{
		VisualStudioPath: vsPath,
	})

	if err != nil {
		return fmt.Errorf("Problem compiling compileCommand template: %s", err)
	}

	argv := []string{fmt.Sprintf("/Fo%s", outfileName)}
	if cvs.Cflags != "" {
		flags, err := shlex.Split(cvs.Cflags)
		if err != nil {
			return fmt.Errorf("Unable to parse cflags: %s", err)
		}

		argv = append(argv, flags...)
	}

	argv = append(argv, srcfileName)
	compileCommand := fmt.Sprintf(`"%s" %s`, buffer.String(), strings.Join(argv, " "))
	envScript := fmt.Sprintf(`"%sVC\bin\vcvars32.bat"`, vsPath)

	compileScript, err := os.Create(filepath.Join(tmpfolder, "compile.bat"))
	if err != nil {
		return fmt.Errorf("Problem creating batch script: %s", err)
	}

	scriptBody := fmt.Sprintf(`call %s
if %%errorlevel%% neq 0 exit /b %%errorlevel%%
call %s
if %%errorlevel%% neq 0 exit /b %%errorlevel%%`, envScript, compileCommand)

	_, err = compileScript.Write([]byte(scriptBody))
	if err != nil {
		return fmt.Errorf("Problem writing batch script: %s", err)
	}

	cmd := exec.Command("cmd.exe", "/c", filepath.Join(tmpfolder, "compile.bat"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Problem running the Visual Studio compile: %s: %s", err, string(out))
	}

	if cvs.Run {
		cmd = exec.Command(outfileName)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Problem running compiled executable: %s: %s", err, string(out))
		}
	}

	return nil
}

// CompileVisualStudioFromArgs will populate the CompileVisualStudio with the args given in the
// tests YAML config
func CompileVisualStudioFromArgs(args Args) (Checker, error) {
	cg := CompileVisualStudio{}

	if err := requiredArgs(args, "source"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &cg); err != nil {
		return nil, err
	}

	return cg, nil
}
