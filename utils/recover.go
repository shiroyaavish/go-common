package utils

import (
	"fmt"
	"runtime/debug"
)

func Recover() {
	if r := recover(); r != nil {
		fmt.Printf("%s", debug.Stack())
	}
}
