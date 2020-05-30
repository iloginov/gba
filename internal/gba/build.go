package gba

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func BuildPackage(pkg string) (string, error) {
	cmd := exec.Command("go", "build", "-a", "-work", pkg)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return "", exiterr
		}
		return "", err
	}
	errStr := string(stderr.Bytes())

	scanner := bufio.NewScanner(strings.NewReader(errStr))
	scanner.Split(bufio.ScanLines)

	var workDir string
	var found bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "WORK=") {
			found = true
			workDir = strings.Split(line, "=")[1]
			break
		}
	}

	if found {
		return workDir, nil
	}

	return "", fmt.Errorf("Can't find build working directory for package %s", pkg)
}
