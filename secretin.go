// Package secretin implements user entry of passwords without echoing.
package secretin

import (
	"fmt"
	"os"
	"syscall"
)

// SetEcho configures user input echo on or off.
func SetEcho(isEcho bool) error {
	var arg string
	if isEcho {
		arg = "echo"
	} else {
		arg = "-echo"
	}

	// These file descriptors need to be passed forward to the forked proc
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	// Despite the fact that the first arg for ForkExec() is the command,
	// the 0th item of the second arg is the command as well, and the 1st
	// item is the first arg for the command.  It seems that the 0th item
	// can be whatever you want.  Perhaps it is the proc name that will
	// show up in ps.
	_, err := syscall.ForkExec("/bin/stty", []string{"stty", arg},
		&syscall.ProcAttr{Files: fds})

	return err
}

// ScanlnNoEcho waits for a line of user input that is not echoed.
func ScanlnNoEcho(a ...interface{}) (int, error) {
	// Turn off stdin echo
	if err := SetEcho(false); err != nil {
		return 0, err
	}

	// Scan in line
	n, err := fmt.Scanln(a...)
	if err != nil {
		return n, err
	}

	// Turn stdin echo back on
	if err := SetEcho(true); err != nil {
		return n, err
	}

	return n, nil
}

// PromptForSecret prompt a user for a line of input that will not be echoed.
func PromptForSecret(msg string) (string, error) {
	fmt.Print(msg)

	var resp string
	if _, err := ScanlnNoEcho(&resp); err != nil {
		return "", err
	}

	// Print new line after prompt since new line wasn't echoed
	fmt.Println()

	return resp, nil
}
