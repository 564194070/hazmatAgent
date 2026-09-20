package tools

import "errors"

type EchoTool struct {
	name string
	desc string
}

func NewEchoTool() *EchoTool {
	return &EchoTool{name: "echo工具", desc: "重复输入的内容，返回给用户"}
}

func (t *EchoTool) Execute(input map[string]any) (string, error) {

	for key, value := range input {
		if key == "input" {
			return value.(string), nil
		}
	}

	return "", errors.New("input 参数不存在，工具调用失败")
}

func (t *EchoTool) GetDesc() string {
	return t.desc
}

func (t *EchoTool) GetName() string {
	return t.name
}
