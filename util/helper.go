package util

import (
	"errors"
	"os"
	"path/filepath"
)

// GetRootDir 获得项目根目录
func GetRootDir() (string, error) {
	// 获得当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上递归查找go.mod文件
	for {
		// 检查当前目录是否存在go.mod
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return dir, nil
		}

		// 到达文件系统根目录仍未找到
		parentDir := filepath.Dir(dir)
		if parentDir == dir { // 已经到达根目录
			return "", errors.New("can not find go.mod file")
		}
		dir = parentDir
	}
}

// GetRootRelPath 获得相对于项目根目录的路径
func GetRootRelPath(relPath string) (string, error) {
	rootDir, err := GetRootDir()

	if err != nil {
		return "", err
	}

	return rootDir + "/" + relPath, nil
}
