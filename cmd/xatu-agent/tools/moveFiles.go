package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type MoveFilesTool struct {
	name string
	desc string
}

func NewMoveFilesTool() *MoveFilesTool {
	return &MoveFilesTool{name: "moveFiles工具", desc: "移动文件"}
}

func (t *MoveFilesTool) Execute(input map[string]any) (string, error) {

	var src, dst string

	if input["src"] != nil {
		src = input["src"].(string)
	} else {
		return "", errors.New("src 参数不存在，工具调用失败")
	}
	if input["dst"] != nil {
		dst = input["dst"].(string)
	} else {
		return "", errors.New("dst 参数不存在，工具调用失败")
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return "", errors.New("src 目录不存在，工具调用失败")
	}

	for _, entry := range entries {
		// 跳过子文件夹，只处理文件
		if entry.IsDir() {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// 移动文件
		err = os.Rename(srcPath, dstPath)
		if err != nil {
			return "", fmt.Errorf("rename %s -> %s failed: %w", src, dst, err)
		}
	}

	return "", nil

}

func (t *MoveFilesTool) GetDesc() string {
	return t.desc
}

func (t *MoveFilesTool) GetName() string {
	return t.name
}
