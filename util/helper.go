package util

import (
	"errors"
	"os"
	"path/filepath"
)

// GetPathFromParent 从逐级父目录获取是否存在目标目录
func GetPathFromParent(targetPath string) (string, error) {
	// 获得当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上递归查找
	for {
		// 检查当前目录是否存在
		nPath := filepath.Join(dir, targetPath)
		if _, err := os.Stat(nPath); err == nil {
			return nPath, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir { // 已经到达根目录
			return "", errors.New("can not find path")
		}
		dir = parentDir
	}
}
