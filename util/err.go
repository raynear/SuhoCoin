package err

import (
    "fmt"
    "runtime"
    "strings"
)

func ERR() {
    pc := make([]uintptr, 10)
    runtime.Callers(1, pc)
    f := runtime.FuncForPC(pc[0])
    fmt.Println("currentFunction:", strings.Split(f.Name(), ".")[1])
    fmt.Println("currentPackage:", strings.Split(f.Name(), ".")[0])
}
