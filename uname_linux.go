package lnksutils

import "syscall"

// UnameMachine return the `uname -m"
func UnameMachine() string {
	var uname syscall.Utsname
	syscall.Uname(&uname)

	arr := uname.Machine[:]
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}
