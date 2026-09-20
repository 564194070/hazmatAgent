package tools

import (
	"errors"
	"fmt"
	"os"
)

type ReadFileListTool struct {
	name string
	desc string
}

func NewReadFileListTool() *ReadFileListTool {
	return &ReadFileListTool{name: "readFileList工具", desc: "读取文件列表"}
}

func (t *ReadFileListTool) Execute(input map[string]any) (string, error) {

	for key, value := range input {
		if key == "path" {

			dir := value.(string)
			entries, err := os.ReadDir(dir)
			if err != nil {
				return "", fmt.Errorf("read dir failed: %w", err)
			}

			var fileNames []string

			for _, entry := range entries {
				// 跳过目录，只保留文件
				if !entry.IsDir() {
					fileNames = append(fileNames, entry.Name())
				}
			}
			return fmt.Sprintf("当前目录%s文件列表: %v", dir, fileNames), nil
		}
	}

	return "", errors.New("path 参数不存在，工具调用失败")
}
