package sessiontest

import (
	"encoding/json"
	"testing"

	opencode "github.com/dominicnunez/opencode-sdk-go"
)

func TestSessionUnmarshal(t *testing.T) {
	raw := `{
		"id": "ses123",
		"directory": "/tmp/test",
		"projectID": "proj456",
		"time": {"created": 1700000000.0},
		"title": "Test Session"
	}`

	var s opencode.Session
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if s.ID != "ses123" {
		t.Errorf("ID = %q, want %q", s.ID, "ses123")
	}
	if s.Directory != "/tmp/test" {
		t.Errorf("Directory = %q, want %q", s.Directory, "/tmp/test")
	}
	if s.ProjectID != "proj456" {
		t.Errorf("ProjectID = %q, want %q", s.ProjectID, "proj456")
	}
	if s.Title != "Test Session" {
		t.Errorf("Title = %q, want %q", s.Title, "Test Session")
	}
	if s.Time.Created != 1700000000.0 {
		t.Errorf("Time.Created = %v, want %v", s.Time.Created, 1700000000.0)
	}
}

func TestAgentPartUnmarshal(t *testing.T) {
	raw := `{
		"id": "part1",
		"messageID": "msg1",
		"name": "test-agent",
		"sessionID": "ses1",
		"type": "agent",
		"source": {"end": 100, "start": 0, "value": "hello"}
	}`

	var ap opencode.AgentPart
	if err := json.Unmarshal([]byte(raw), &ap); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if ap.ID != "part1" {
		t.Errorf("ID = %q, want %q", ap.ID, "part1")
	}
	if ap.Name != "test-agent" {
		t.Errorf("Name = %q, want %q", ap.Name, "test-agent")
	}
	if ap.Type != opencode.AgentPartTypeAgent {
		t.Errorf("Type = %q, want %q", ap.Type, opencode.AgentPartTypeAgent)
	}
	if ap.Source.End != 100 {
		t.Errorf("Source.End = %d, want %d", ap.Source.End, 100)
	}
	if ap.Source.Start != 0 {
		t.Errorf("Source.Start = %d, want %d", ap.Source.Start, 0)
	}
	if ap.Source.Value != "hello" {
		t.Errorf("Source.Value = %q, want %q", ap.Source.Value, "hello")
	}
}

func TestAssistantMessageUnmarshal(t *testing.T) {
	raw := `{
		"id": "msg1",
		"cost": 0.01,
		"mode": "chat",
		"modelID": "gpt-4",
		"parentID": "parent1",
		"path": {"cwd": "/home", "root": "/"},
		"providerID": "openai",
		"role": "assistant",
		"sessionID": "ses1",
		"system": ["sys1"],
		"time": {"created": 1000.0, "completed": 2000.0},
		"tokens": {
			"cache": {"read": 10, "write": 20},
			"input": 100,
			"output": 50,
			"reasoning": 30
		},
		"summary": false
	}`

	var am opencode.AssistantMessage
	if err := json.Unmarshal([]byte(raw), &am); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if am.ID != "msg1" {
		t.Errorf("ID = %q, want %q", am.ID, "msg1")
	}
	if am.ModelID != "gpt-4" {
		t.Errorf("ModelID = %q, want %q", am.ModelID, "gpt-4")
	}
	if am.Path.Cwd != "/home" {
		t.Errorf("Path.Cwd = %q, want %q", am.Path.Cwd, "/home")
	}
	if am.Role != opencode.AssistantMessageRoleAssistant {
		t.Errorf("Role = %q, want %q", am.Role, opencode.AssistantMessageRoleAssistant)
	}
	if am.Time.Created != 1000.0 {
		t.Errorf("Time.Created = %v, want %v", am.Time.Created, 1000.0)
	}
	if am.Time.Completed != 2000.0 {
		t.Errorf("Time.Completed = %v, want %v", am.Time.Completed, 2000.0)
	}
	if am.Tokens.Input != 100 {
		t.Errorf("Tokens.Input = %v, want %v", am.Tokens.Input, 100.0)
	}
	if am.Tokens.Cache.Read != 10 {
		t.Errorf("Tokens.Cache.Read = %v, want %v", am.Tokens.Cache.Read, 10.0)
	}
	if len(am.System) != 1 || am.System[0] != "sys1" {
		t.Errorf("System = %v, want [sys1]", am.System)
	}
}

func TestStepFinishPartUnmarshal(t *testing.T) {
	raw := `{
		"id": "step1",
		"cost": 0.05,
		"messageID": "msg1",
		"reason": "done",
		"sessionID": "ses1",
		"tokens": {
			"cache": {"read": 5, "write": 10},
			"input": 200,
			"output": 100,
			"reasoning": 50
		},
		"type": "step-finish",
		"snapshot": "snap123"
	}`

	var sf opencode.StepFinishPart
	if err := json.Unmarshal([]byte(raw), &sf); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if sf.ID != "step1" {
		t.Errorf("ID = %q, want %q", sf.ID, "step1")
	}
	if sf.Cost != 0.05 {
		t.Errorf("Cost = %v, want %v", sf.Cost, 0.05)
	}
	if sf.Reason != "done" {
		t.Errorf("Reason = %q, want %q", sf.Reason, "done")
	}
	if sf.Type != opencode.StepFinishPartTypeStepFinish {
		t.Errorf("Type = %q, want %q", sf.Type, opencode.StepFinishPartTypeStepFinish)
	}
	if sf.Snapshot != "snap123" {
		t.Errorf("Snapshot = %q, want %q", sf.Snapshot, "snap123")
	}
	if sf.Tokens.Input != 200 {
		t.Errorf("Tokens.Input = %v, want %v", sf.Tokens.Input, 200.0)
	}
}

func TestUserMessageUnmarshal(t *testing.T) {
	raw := `{
		"id": "umsg1",
		"role": "user",
		"sessionID": "ses1",
		"time": {"created": 1500.0}
	}`

	var um opencode.UserMessage
	if err := json.Unmarshal([]byte(raw), &um); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if um.ID != "umsg1" {
		t.Errorf("ID = %q, want %q", um.ID, "umsg1")
	}
	if um.Role != opencode.UserMessageRoleUser {
		t.Errorf("Role = %q, want %q", um.Role, opencode.UserMessageRoleUser)
	}
	if um.Time.Created != 1500.0 {
		t.Errorf("Time.Created = %v, want %v", um.Time.Created, 1500.0)
	}
}

func TestToolPartUnmarshal(t *testing.T) {
	raw := `{
		"id": "tool1",
		"messageID": "msg1",
		"sessionID": "ses1",
		"tool": "bash",
		"type": "tool",
		"input": {"command": "ls"},
		"metadata": {"key": "val"}
	}`

	var tp opencode.ToolPart
	if err := json.Unmarshal([]byte(raw), &tp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if tp.ID != "tool1" {
		t.Errorf("ID = %q, want %q", tp.ID, "tool1")
	}
	if tp.Tool != "bash" {
		t.Errorf("Tool = %q, want %q", tp.Tool, "bash")
	}
	if tp.Type != opencode.ToolPartTypeTool {
		t.Errorf("Type = %q, want %q", tp.Type, opencode.ToolPartTypeTool)
	}
}

func TestSessionCommandResponseUnmarshal(t *testing.T) {
	raw := `{
		"info": {
			"id": "msg1",
			"model": "gpt-4",
			"path": {"cwd": "/home", "root": "/"},
			"role": "assistant",
			"sessionID": "ses1",
			"system": [],
			"time": {"created": 1000.0},
			"tokens": {
				"cache": {"read": 0, "write": 0},
				"input": 0,
				"output": 0,
				"reasoning": 0
			}
		},
		"parts": []
	}`

	var resp opencode.SessionCommandResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Info.ID != "msg1" {
		t.Errorf("Info.ID = %q, want %q", resp.Info.ID, "msg1")
	}
	if len(resp.Parts) != 0 {
		t.Errorf("Parts length = %d, want 0", len(resp.Parts))
	}
}

func TestSessionUnmarshal_OptionalFields(t *testing.T) {
	// Test with only required fields — optional fields should be zero values
	raw := `{
		"id": "ses1",
		"directory": "/tmp",
		"projectID": "proj1",
		"time": {"created": 100.0},
		"title": "minimal"
	}`

	var s opencode.Session
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if s.ID != "ses1" {
		t.Errorf("ID = %q, want %q", s.ID, "ses1")
	}
	// ParentID should be zero value when not present
	if s.ParentID != "" {
		t.Errorf("ParentID = %q, want empty string", s.ParentID)
	}
}

func TestReasoningPartUnmarshal(t *testing.T) {
	raw := `{
		"id": "reas1",
		"messageID": "msg1",
		"sessionID": "ses1",
		"type": "reasoning",
		"time": {"start": 100.0, "end": 200.0},
		"metadata": {"key": "value"}
	}`

	var rp opencode.ReasoningPart
	if err := json.Unmarshal([]byte(raw), &rp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if rp.ID != "reas1" {
		t.Errorf("ID = %q, want %q", rp.ID, "reas1")
	}
	if rp.Type != opencode.ReasoningPartTypeReasoning {
		t.Errorf("Type = %q, want %q", rp.Type, opencode.ReasoningPartTypeReasoning)
	}
	if rp.Time.Start != 100.0 {
		t.Errorf("Time.Start = %v, want %v", rp.Time.Start, 100.0)
	}
	if rp.Time.End != 200.0 {
		t.Errorf("Time.End = %v, want %v", rp.Time.End, 200.0)
	}
}
