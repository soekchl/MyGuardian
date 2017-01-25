// MyAutoGuardian project main.go

// For Windows

package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	. "github.com/soekchl/myUtils"
)

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

func main() {
	SetOutputFileLog("MyAutoGuardian", true)

	if len(os.Args) == 2 {
		Notice(os.Args)
	} else {
		Error(os.Args[0], " [ProcessName]")
		return
	}

	result, err := IsHave(os.Args[1])
	n := 0

	for ; ; time.Sleep(time.Second) {
		result, err = IsHave(os.Args[1])
		if err != nil {
			Error(err)
			continue
		}
		if !result {
			Warn(os.Args[1], " Restarting!!! time:", n)
			n++

			cmd := exec.Command(os.Args[1])
			err = cmd.Run()
			if err != nil {
				Error(err)
				continue
			}
		}
	}
	Notice("My Jobs Is Over!!!")
}

func IsHave(key string) (result bool, err error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		err = errors.New("int(pHandle) == -1")
		return
	}
	Process32Next := kernel32.NewProc("Process32Next")
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			if strings.Contains(string(proc.szExeFile[0:]), key) {
				result = true
				break
			}
		} else {
			break
		}
	}

	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(pHandle)
	return
}
