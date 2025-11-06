package warp

import (
	"os/exec"
	"runtime"
)

func ping(addr string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows ping default uses -n for count
		cmd = exec.Command("ping", "-n", "1", addr)
	} else {
		// Linux/macOS use -c for count
		cmd = exec.Command("ping", "-c", "1", addr)
	}
	debug.Println("ping", addr)
	if err := cmd.Run(); err != nil {
		debug.Print(err)
	}
}
