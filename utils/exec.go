package utils

import (
	"bytes"
	"os/exec"
)

// RunShellScript executes a shell script with arguments and returns the output as a string and any errors encountered.
// The function takes two arguments, `script` which is the shell script to be executed and `args` which is a variadic argument of type string representing script arguments.
// It returns two values, a string containing the output of the script and an error object. If the execution of the script encounters any error, the string will be empty and the error will be non-nil
func RunShellScript(script string, args ...string) (string, error) {
	// Append the script and the args
	args = append([]string{"-c", script}, args...)

	cmd := exec.Command("/bin/bash", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
