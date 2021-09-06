package tools

import (
	"runtime"
	"strings"
)

//FuncName ..
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcNameEnd := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return funcNameEnd[len(funcNameEnd)-1] + "():"
}
