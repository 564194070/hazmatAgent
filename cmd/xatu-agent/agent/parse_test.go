package agent

import "testing"

func TestParseReActAction(t *testing.T) {
	llmResp := "Thought: 先列出目录。\nAction: readFileList[/root/agent/frame]"
	action, err := parseReActAction(llmResp)
	if err != nil {
		t.Fatal(err)
	}
	if action.Name != "readFileList" {
		t.Fatalf("Name = %q", action.Name)
	}
	if action.Input != "/root/agent/frame" {
		t.Fatalf("Input = %q", action.Input)
	}
	if action.IsFinish() {
		t.Fatal("should not be Finish")
	}
}

func TestParseReActActionFinish(t *testing.T) {
	llmResp := "Thought: 完成。\nAction: Finish[已创建全部目录]"
	action, err := parseReActAction(llmResp)
	if err != nil {
		t.Fatal(err)
	}
	if !action.IsFinish() {
		t.Fatal("should be Finish")
	}
	if action.Input != "已创建全部目录" {
		t.Fatalf("Input = %q", action.Input)
	}
}
