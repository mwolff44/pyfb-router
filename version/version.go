package version

import (
	"fmt"
	"runtime"
)

// GitCommit is The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version is The main version number that is being run at the moment.
const Version = "0.2.0"

// BuildDate is the build date
var BuildDate = ""

// GoVersion is the go version
var GoVersion = runtime.Version()

//OsArch is the OS architecture
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
