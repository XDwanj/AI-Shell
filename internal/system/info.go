package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// OsInfo 存储系统信息
type OsInfo struct {
	DistroID      string
	DistroVersion string
	PrettyName    string
}

// GetOsInfo 获取系统信息
func GetOsInfo() (*OsInfo, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("读取系统信息失败: %v", err)
	}
	defer file.Close()

	info := &OsInfo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := strings.Trim(parts[1], "\"")

		switch key {
		case "ID":
			info.DistroID = value
		case "VERSION_ID":
			info.DistroVersion = value
		case "PRETTY_NAME":
			info.PrettyName = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析系统信息失败: %v", err)
	}

	return info, nil
}

// GetUserID 获取用户ID信息
func GetUserID() (string, error) {
	cmd := exec.Command("id")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取用户ID失败: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDirectoryInfo 获取目录信息
func GetDirectoryInfo() (pwd string, files string, err error) {
	// 获取当前工作目录
	pwd, err = os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 获取目录内容
	entries, err := os.ReadDir(".")
	if err != nil {
		return pwd, "", fmt.Errorf("读取目录内容失败: %v", err)
	}

	var filesList []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		filesList = append(filesList, name)
	}

	files = strings.Join(filesList, " ")
	return pwd, files, nil
}

// GetSystemInfo 获取完整的系统信息
func GetSystemInfo() (string, error) {
	osInfo, err := GetOsInfo()
	if err != nil {
		return "", err
	}

	userID, err := GetUserID()
	if err != nil {
		return "", err
	}

	pwd, files, err := GetDirectoryInfo()
	if err != nil {
		return "", err
	}

	shell := os.Getenv("SHELL")

	info := fmt.Sprintf(`[echo $SHELL]
%s
[系统信息]
发行版: %s
发行版ID: %s
版本: %s
[id]
%s
[ls -aF]
%s
[pwd]
%s`,
		shell,
		osInfo.PrettyName,
		osInfo.DistroID,
		osInfo.DistroVersion,
		userID,
		files,
		pwd)

	return info, nil
}
