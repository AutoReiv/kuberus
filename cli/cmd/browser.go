package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	fmt.Println("Executing command:", cmd, args)
	return exec.Command(cmd, args...).Start()
}
