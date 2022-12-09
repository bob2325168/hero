package util

import (
	"log"
	"os"
	"syscall"
)

// GetExecDirectory 获取当前执行程序的路径
func GetExecDirectory() string {
	file, err := os.Getwd()
	log.Println("本地项目路径: ", file, err)
	if err == nil {
		return file + "/"
	}
	return ""
}

// CheckProcessExist 检测进程的PID是否存在
func CheckProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}
