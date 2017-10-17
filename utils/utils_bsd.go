// +build bsd darwin

package utils

import (
	"syscall"
	"unsafe"
)

// IsaTty tells whether given fd is a terminal
func IsaTty(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
