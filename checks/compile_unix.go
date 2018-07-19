// +build !windows

package checks

func init() {
	availableChecks["compile-gcc"] = func(args Args) (Checker, error) {
		if _, provided := args["compiler"]; !provided {
			args["compiler"] = "gcc"
		}

		return CompileChecker{}.FromArgs(args)
	}
}
