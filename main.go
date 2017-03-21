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

// 获取 Windows进程列表结构体
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
	// 设定输出 Log 名字               是否文件输出
	SetOutputFileLog("MyAutoGuardian", true)

	// 判断后台参数是否为 2个
	if len(os.Args) == 2 {
		Notice(os.Args)
	} else {
		Error(os.Args[0], " [ProcessName]")
		return
	}

	//　判断守护进程是否挂掉
	result, err := IsHave(os.Args[1])
	n := 0
	var sleep_time time.Duration = 1

	// 每秒监控一次
	for ; ; time.Sleep(time.Second * sleep_time) {
		result, err = IsHave(os.Args[1])
		if err != nil {
			Error(err)
			n = 0
			continue
		}
		// 进程挂掉
		if !result {
			Warn(os.Args[1], " Restarting!!! times:", n)
			n++

			// 重启
			cmd := exec.Command(os.Args[1])
			err = cmd.Run()
			if err != nil {
				Error(err)
				continue
			}
		}

		if n > 10 {
			sleep_time = time.Duration(n - 10)
			if sleep_time > 5 {
				sleep_time = 5
			}
		}
	}
	Notice("My Jobs Is Over!!!")
}

// 调用WinAPI 获取进程列表
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
