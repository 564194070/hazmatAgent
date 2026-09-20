package tools

import (
	"errors"
	"os"
)

type MkdirAllTool struct {
	name string
	desc string
}

func NewMkdirAllTool() *MkdirAllTool {
	return &MkdirAllTool{name: "mkdirAll工具", desc: "创建目录"}
}

func (t *MkdirAllTool) Execute(input map[string]any) (string, error) {

	for key, value := range input {
		if key == "path" {

			path := value.(string)
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return "", err
			}
			return "", nil
		}
	}

	return "", errors.New("path 参数不存在，工具调用失败")
}

func (t *MkdirAllTool) GetDesc() string {
	return t.desc
}

func (t *MkdirAllTool) GetName() string {
	return t.name
}
