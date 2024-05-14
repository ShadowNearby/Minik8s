package utils

import (
	"fmt"
	"path"
	"runtime"
)

var _, _filename, _, _ = runtime.Caller(0)

var RootPath = path.Dir(path.Dir(_filename))

var ConfigPath = fmt.Sprintf("%s/%s", RootPath, "config")

var ExamplePath = fmt.Sprintf("%s/%s", RootPath, "example")
