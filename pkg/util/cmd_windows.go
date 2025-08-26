//go:build windows

package util

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
	MAX_PATH           = 260
)

type ProcessEntry32 struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [MAX_PATH]uint16 // Use uint16 instead of byte
}

func RunCommand(ctx context.Context, text string) error {
	fmt.Println(text)
	var cmd *exec.Cmd
	if os.PathSeparator == '\\' {
		cmd = exec.CommandContext(ctx, "cmd", "/c", text)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", text)
	}
	// Capture standard output
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	// Execute command cmd.CombinedOutput() and capture output
	//output, err = cmd.CombinedOutput()
	if err = cmd.Start(); err != nil {
		return err
	}
	readout := bufio.NewReader(stdout)
	GetOutput(readout)
	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func GetOutput(reader *bufio.Reader) {
	var sumOutput string // Collect all screen output content
	outputBytes := make([]byte, 200)
	for {
		n, err := reader.Read(outputBytes) // Get real-time screen output (not split by newlines, so combine with sumOutput)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			sumOutput += err.Error()
		}
		output := string(outputBytes[:n])
		fmt.Print(output) // Output screen content
		sumOutput += output
	}
	return
}

// PrintSleepTime prints 0-60 second wait
func PrintSleepTime(sec int) {
	if sec <= 0 || sec > 60 {
		return
	}
	fmt.Println()
	for t := sec; t > 0; t-- {
		seconds := strconv.Itoa(int(t))
		if t < 10 {
			seconds = fmt.Sprintf("0%d", t)
		}
		fmt.Printf("\rplease wait.... [00:%s of appr. Max %d sec]", seconds, sec)
		time.Sleep(time.Second)
	}
	fmt.Println()
}

func StartProcess(inputUri string, outfile string, args []string) bool {
	//	procAttr := &os.ProcAttr{
	//		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	//	}
	//	userArgs := strings.Split(config.Conf.DezoomifyRs, " ")
	//	argv := []string{"/c", config.Conf.DezoomifyPath}
	//	if userArgs != nil {
	//		argv = append(argv, userArgs...)
	//	}
	//	if args != nil {
	//		argv = append(argv, args...)
	//	}
	//	argv = append(argv, inputUri, outfile)
	//	process, err := os.StartProcess("C:\\Windows\\System32\\cmd.exe", argv, procAttr)
	//	if err != nil {
	//		fmt.Println("start process error:", err)
	//		return false
	//	}
	//	_, err = process.Wait()
	//	if err != nil {
	//		fmt.Println("wait error:", err)
	//		return false
	//	}
	//	fmt.Println()
	return false
}

func OpenWebBrowser(args []string) bool {
	running, err := IsBookgetGuiRunning()
	if running || err != nil {
		return false
	}
	go func() {
		procAttr := &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		}
		fPath, _ := os.Executable()
		binDir := filepath.Dir(fPath)

		argv := []string{"/c", "-i"}
		if args != nil {
			argv = append(argv, args...)
		}
		process, err := os.StartProcess(binDir+"\\bookget-gui.exe", argv, procAttr)
		if err != nil {
			fmt.Println("start process error:", err)
			return
		}
		_ = process.Release()
	}()

	return true
}

func isProcessRunning(processName string) (bool, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	createToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	process32First := kernel32.NewProc("Process32FirstW") // 使用宽字符版本
	process32Next := kernel32.NewProc("Process32NextW")   // 使用宽字符版本
	closeHandle := kernel32.NewProc("CloseHandle")

	// 创建进程快照
	snapshot, _, err := createToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == uintptr(syscall.InvalidHandle) {
		return false, fmt.Errorf("failed to create process snapshot: %v", err)
	}
	defer closeHandle.Call(snapshot)

	var entry ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	// 获取第一个进程
	ret, _, err := process32First.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return false, fmt.Errorf("failed to get process info: %v", err)
	}

	for {
		// Correctly convert UTF-16 string
		exeFile := syscall.UTF16ToString(entry.ExeFile[:])
		// Remove null characters at end of string
		exeFile = strings.TrimRight(exeFile, "\x00")
		if strings.EqualFold(exeFile, processName) {
			return true, nil
		}

		// Get next process
		ret, _, err := process32Next.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			if errno, ok := err.(syscall.Errno); ok && errno == 0 {
				break // Normal end
			}
			return false, nil
		}
	}

	return false, nil
}

func IsBookgetGuiRunning() (ok bool, err error) {
	processName := "bookget-gui.exe"
	running, err := isProcessRunning(processName)
	if err != nil {
		return false, err
	}

	if running {
		return true, nil
	} else {
		return false, nil
	}
}
