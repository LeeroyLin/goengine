package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PathExists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	// 没报错，则存在
	if err == nil {
		return true, nil
	}

	// 报错路径不存在
	if errors.Is(err, os.ErrNotExist) {
		return false, errors.New("path not exist. " + path)
	}

	return false, err
}

// DirExists 判断目录是否存在
func DirExists(direction string) (bool, error) {
	info, err := os.Stat(direction)
	if err != nil {
		return false, err
	}

	// 确保是目录而不是文件
	return info.IsDir(), nil
}

// IsEmpty 字符串是否是空字符串
func IsEmpty(s string) bool {
	return s == ""
}

// IsEmptyOrWhitespace 字符串是否是空字符串或只包含空格
func IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Mkdir 创建目录
func Mkdir(dir string) error {
	// 目录存在
	if ok, _ := DirExists(dir); ok {
		return nil
	}

	if err := os.MkdirAll(dir, 0775); err != nil {
		if os.IsPermission(err) {
			return err
		}
	}

	return nil
}

// IntToFixedStr 将int类型转换为指定长度的字符串
func IntToFixedStr(val, length int) string {
	varStr := strconv.Itoa(val)
	lenVal := len(varStr)

	if length > 0 {
		if lenVal < length {
			return fmt.Sprintf("%0"+strconv.Itoa(length)+"d", val)
		} else {
			return varStr[:length]
		}
	}

	return strconv.Itoa(val)
}
