package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// OpenURL opens the specified URL in the default browser
func OpenURL(url string) error {
	// Skip browser open in test mode
	if os.Getenv("NEXLAYER_TEST_MODE") == "true" {
		return nil
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
