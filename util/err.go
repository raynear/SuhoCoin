package util

import (
	"fmt"
	"runtime"
	"strings"
)

func ERR(msg string, e error) {
	if e != nil {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		pwd := strings.Split(f.Name(), ".")
		fmt.Printf("!!ERR!! ")
		fmt.Printf("Pack:%s", pwd[0])
		if len(pwd) > 2 {
			fmt.Printf(" Func:%s.%s | ", pwd[1], pwd[2])
			fmt.Println(msg)
		} else {
			fmt.Println(" Func:", pwd[1], "|", msg)
		}
		fmt.Println(e)
	}
}
