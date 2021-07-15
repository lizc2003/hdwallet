package qa

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func IsProcessRunning(port uint, name string) bool {
	if strings.Contains(runtime.GOOS, "linux") {
		checkPortCmd := exec.Command("netstat", "-ntpl")

		cmdPrint := executeCmd(checkPortCmd)
		if strings.Contains(cmdPrint, strconv.Itoa(int(port))) && strings.Contains(cmdPrint, name) {
			return true
		}
		return false
	} else if strings.Contains(runtime.GOOS, "darwin") {
		checkPortCmd := exec.Command("lsof", "-i", "tcp:"+strconv.FormatInt(int64(port), 10))
		cmdPrint := executeCmd(checkPortCmd)
		if strings.Contains(cmdPrint, strconv.Itoa(int(port))) && strings.Contains(cmdPrint, "bitcoin") {
			return true
		}
		return false
	} else {
		panic("unsupported platform")
	}
}

func executeCmd(cmd *exec.Cmd) string {
	stderr, err := cmd.StderrPipe()
	panicIf(err, "Failed to get stderr pip ")
	stdout, err := cmd.StdoutPipe()
	panicIf(err, fmt.Sprintf("Failed to get stdout pipe %v", err))

	err = cmd.Start()
	panicIf(err, fmt.Sprintf("Failed to start cmd %v", err))

	b, err := ioutil.ReadAll(stdout)
	panicIf(err, fmt.Sprintf("Failed to read cmd (%v) stdout, %v", cmd, err))
	out := string(b)

	bo, err := ioutil.ReadAll(stderr)
	panicIf(err, "Failed to read stderr")
	out += string(bo)

	cmd.Wait()
	stdout.Close()
	stderr.Close()
	return strings.TrimSpace(out)
}

func panicIf(e error, msg string) {
	if e != nil {
		panic(fmt.Errorf("【ERR】 %s %v", msg, e))
	}
}
