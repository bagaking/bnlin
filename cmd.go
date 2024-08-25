package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/bagaking/botheater/utils"
	"github.com/khicago/irr"
)

const (
	VersionFailed  = "获取版本信息失败"
	VersionError   = "无法解析版本信息"
	VersionUnknown = "未知版本"
)

// getOSInfo 获取操作系统信息和命令行语言
func getOSInfo() (osType, version, lang string) {
	osType = runtime.GOOS
	switch osType {
	case "windows":
		cmd := "ver"
		out, err := exec.Command(cmd).Output()
		if err != nil {
			version = VersionFailed
		} else {
			version = strings.TrimSpace(string(out))
		}
	case "linux":
		cmd := "cat /proc/version"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			version = VersionFailed
		} else {
			parts := strings.Fields(string(out))
			if len(parts) >= 3 {
				version = parts[2]
			} else {
				version = VersionError
			}
		}
	case "darwin":
		cmd := "sw_vers"
		out, err := exec.Command(cmd).Output()
		if err != nil {
			version = VersionFailed
		} else {
			lines := bytes.Split(out, []byte("\n"))
			for _, line := range lines {
				if strings.Contains(string(line), "ProductVersion:") {
					version = strings.TrimSpace(strings.TrimPrefix(string(line), "ProductVersion: "))
					break
				}
			}
		}
		// 获取系统语言
		langCmd := "defaults read -g AppleLocale"
		langOut, err := exec.Command("bash", "-c", langCmd).Output()
		if err != nil {
			lang = "unknown"
		} else {
			lang = strings.TrimSpace(string(langOut))
		}
	default:
		version = VersionUnknown
	}

	// 获取命令行语言
	if lang == "" {
		lang = os.Getenv("LANG")
		if lang == "" {
			lang = os.Getenv("LC_ALL")
		}
		if lang == "" {
			lang = "unknown"
		}
	}

	return osType, version, lang
}

func execute(comment string) error {
	// Extract the actual command from the comment
	lines := strings.Split(comment, "\n")
	var scriptContent strings.Builder
	var commentContent strings.Builder
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || l == "#" || l == "```" || l == "```bash" || l == "```sh" {
			continue
		}
		if strings.HasPrefix(l, "#") {
			commentContent.WriteString(line + "\n")
			continue
		}

		scriptContent.WriteString(line + "\n")
	}

	fmt.Println(utils.SPrintWithCallStack("execution plan", strings.TrimSpace(commentContent.String()), 180))

	// Write the script content to a temporary file
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		return irr.Wrap(err, "failed to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(scriptContent.String()); err != nil {
		return irr.Wrap(err, "failed to write to temp file")
	}
	if err = tmpFile.Close(); err != nil {
		return irr.Wrap(err, "failed to close temp file")
	}

	// Make the script executable
	if err = os.Chmod(tmpFile.Name(), 0o755); err != nil {
		return irr.Wrap(err, "failed to make temp file executable")
	}

	// Execute the script
	cmd := exec.Command("bash", tmpFile.Name())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return irr.Wrap(err, "failed to get stdout pipe")
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return irr.Wrap(err, "failed to get stderr pipe")
	}

	if err = cmd.Start(); err != nil {
		return irr.Wrap(err, "failed to start command")
	}

	// Stream the output
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderrPipe)
	}()

	// Wait for the command to finish
	if err = cmd.Wait(); err != nil {
		return irr.Wrap(err, "command execution failed")
	}

	// Wait for all output to be copied
	wg.Wait()

	return nil
}
