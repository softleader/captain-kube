package command

import "os/exec"

// IsAvailable 檢查是否有傳入的 name command
func IsAvailable(name string) bool {
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
