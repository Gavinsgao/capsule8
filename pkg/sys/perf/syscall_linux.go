// Copyright 2017 Capsule8, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package perf

import (
	"bytes"
	"unsafe"

	"golang.org/x/sys/unix"
)

func enable(fd int) error {
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), PERF_EVENT_IOC_ENABLE, 1); errno != 0 {
		return errno
	}
	return nil
}

func setFilter(fd int, filter string) error {
	if f, err := unix.BytePtrFromString(filter); err != nil {
		return err
	} else if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), PERF_EVENT_IOC_SET_FILTER, uintptr(unsafe.Pointer(f))); errno != 0 {
		return errno
	}
	return nil
}

func disable(fd int) error {
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), PERF_EVENT_IOC_DISABLE, 1); errno != 0 {
		return errno
	}
	return nil
}

func open(attr *EventAttr, pid int, cpu int, groupFd int, flags uintptr) (int, error) {
	buf := new(bytes.Buffer)
	attr.write(buf)
	b := buf.Bytes()

	r1, _, errno := unix.Syscall6(unix.SYS_PERF_EVENT_OPEN, uintptr(unsafe.Pointer(&b[0])),
		uintptr(pid), uintptr(cpu), uintptr(groupFd), uintptr(flags), uintptr(0))
	if errno != 0 {
		return -1, errno
	}
	return int(r1), nil
}
