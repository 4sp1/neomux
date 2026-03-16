//go:build darwin

package adapter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/4sp1/neomux/internal/domain/proc"
	"github.com/4sp1/neomux/internal/repo"
)

type adapter struct{}

func New() (repo.Proc, error) {
	return &adapter{}, nil
}

func (a adapter) List() ([]proc.Proc, error) {
	return processes()
}

// https://github.com/mitchellh/go-ps/blob/ddafa7589c607e0e81a7436520ccdbac913665a2/process_darwin.go#L45
func processes() ([]proc.Proc, error) {
	buf, err := kernProcAll()
	if err != nil {
		return nil, fmt.Errorf("kern proc all: %w", err)
	}
	procs := make([]*kinfoProc, 0, 50)
	k := 0
	for i := _KINFO_STRUCT_SIZE; i < buf.Len(); i += _KINFO_STRUCT_SIZE {
		proc := &kinfoProc{}
		err = binary.Read(bytes.NewBuffer(buf.Bytes()[k:i]), binary.LittleEndian, proc)
		if err != nil {
			return nil, err
		}
		k = i
		procs = append(procs, proc)
	}
	darwinProcs := make([]proc.Proc, len(procs))
	for i, p := range procs {
		darwinProcs[i] = proc.Proc{
			PID:    int(p.Pid),
			Binary: darwinCstring(p.Comm),
		}
	}
	return darwinProcs, nil
}

// https://github.com/mitchellh/go-ps/blob/ddafa7589c607e0e81a7436520ccdbac913665a2/process_darwin.go#L130
type kinfoProc struct {
	_    [40]byte
	Pid  int32
	_    [199]byte
	Comm [16]byte
	_    [301]byte
	PPid int32
	_    [84]byte
}

// https://github.com/mitchellh/go-ps/blob/ddafa7589c607e0e81a7436520ccdbac913665a2/process_darwin.go#L76
func darwinCstring(s [16]byte) string {
	i := 0
	for _, b := range s {
		if b != 0 {
			i++
		} else {
			break
		}
	}

	return string(s[:i])
}

// https://github.com/mitchellh/go-ps/blob/ddafa7589c607e0e81a7436520ccdbac913665a2/process_darwin.go#L89
func kernProcAll() (*bytes.Buffer, error) {
	mib := [4]int32{_CTRL_KERN, _KERN_PROC, _KERN_PROC_ALL, 0}
	size := uintptr(0)

	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		4,
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		0)
	if errno != 0 {
		return nil, errno
	}

	bs := make([]byte, size)
	_, _, errno = syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		4,
		uintptr(unsafe.Pointer(&bs[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		0)
	if errno != 0 {
		return nil, errno
	}

	return bytes.NewBuffer(bs[:size]), nil
}

const (
	_CTRL_KERN         = 1
	_KERN_PROC         = 14
	_KERN_PROC_ALL     = 0
	_KINFO_STRUCT_SIZE = 648
)
