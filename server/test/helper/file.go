package helper

import (
	"fmt"
	"os"
	"path"
)

// 复制文件
func CopyFile(src string, dst string, fileName string) error {
	// 读取源文件
	srcFile := path.Join(src, fileName)
	data, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("读取文件失败：%v", err)
	}

	// 写入目标文件
	dstFile := path.Join(dst, fileName)
	err = os.WriteFile(dstFile, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败：%v", err)
	}
	return nil
}
