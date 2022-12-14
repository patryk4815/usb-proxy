package rawgadget

import (
	"reflect"
	"syscall"
	"unsafe"
)

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRMASK   = ((1 << _IOC_NRBITS) - 1)
	_IOC_TYPEMASK = ((1 << _IOC_TYPEBITS) - 1)
	_IOC_SIZEMASK = ((1 << _IOC_SIZEBITS) - 1)
	_IOC_DIRMASK  = ((1 << _IOC_DIRBITS) - 1)

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = (_IOC_NRSHIFT + _IOC_NRBITS)
	_IOC_SIZESHIFT = (_IOC_TYPESHIFT + _IOC_TYPEBITS)
	_IOC_DIRSHIFT  = (_IOC_SIZESHIFT + _IOC_SIZEBITS)

	_IOC_NONE  = 0
	_IOC_WRITE = 1
	_IOC_READ  = 2
)

func _IOC(dir int, t int, nr int, size int) int {
	return (dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT)
}

func _IO(t int, nr int) int {
	return _IOC(_IOC_NONE, t, nr, 0)
}

func Sizeof(stru interface{}) int {
	c := reflect.TypeOf(stru)
	return int(c.Size())
}

func _IOR(t int, nr int, size int) int {
	return _IOC(_IOC_READ, t, nr, size)
}

func _IOW(t int, nr int, size int) int {
	return _IOC(_IOC_WRITE, t, nr, size)
}

func _IOWR(t int, nr int, size int) int {
	return _IOC(_IOC_WRITE|_IOC_READ, t, nr, size)
}

func ioctlPtr(fd int, req int, ptr unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(ptr))
}

func ioctlInt(fd int, req int, ptr uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(ptr))
}
